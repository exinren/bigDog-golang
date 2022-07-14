package common

import (
	"bigDog-golang/pkg/setting"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"time"
)

type PolygonClient struct {
	NodeURL string
	Client *ethclient.Client
	KeyPath string
}

type Wallet struct {
	KS *keystore.KeyStore	// 账号管理
	Address *accounts.Account		// 公钥地址
	Key string		// 私钥
}

// 初始化客户端
// nodeURL 测试节点地址
func InitPolygonClient(polygon *setting.PolygonSettingS) (*PolygonClient,error){
	client, err := ethclient.Dial(polygon.NodeURL)
	if nil != err {
		return nil,err
	}
	instance := &PolygonClient{
		Client: client,
		NodeURL: polygon.NodeURL,
		KeyPath: polygon.KeyStorePath,
	}
	return instance,nil
}

// CreateWallet使用完，一定 defer CloseWallet(wallet)，否则下次导入的时候包会存在
func (instance *PolygonClient)  CreateWallet() (*Wallet, error)  {
	ks := keystore.NewKeyStore(instance.KeyPath, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount("")
	// 转成json格式
	keyJSON, err := ks.Export(account, "", "")
	// 把JSON格式的转成privateKey
	key, err := keystore.DecryptKey(keyJSON, "")
	if err != nil {
		return nil, err
	}
	//	获取私钥
	privateKey := hex.EncodeToString(crypto.FromECDSA(key.PrivateKey))
	wallet := &Wallet{ks,&account,privateKey}
	return wallet, nil
}

//	ImportWallet 导入钱包，导入后一定 defer CloseWallet(wallet)
//	key	私钥
//	*Wallet	钱包
//	error 不为nil时则钱包导入失败
func (instance *PolygonClient) ImportWallet(key string) (*Wallet,error){
	ks := keystore.NewKeyStore(instance.KeyPath, keystore.StandardScryptN, keystore.StandardScryptP)
	// 私钥转成privateKey结构体
	privateKey,err := crypto.HexToECDSA(key)
	account, err := ks.ImportECDSA(privateKey,"")

	if err != nil {
		return nil,err
	}
	// 延迟解锁
	ks.TimedUnlock(account,"",1000 * time.Second)

	// 返回钱包结构体
	wallet := &Wallet{ks,&account,key}
	return wallet,nil
}

//	CloseWallet 删除本地私钥文件
func (instance *PolygonClient) CloseWallet(wallet *Wallet){
	if nil != wallet {
		wallet.KS.Delete(*wallet.Address, "")
	}
}
// 合约地址：0xe9f4c1CB0F71B4c59612Dab5A0C0Ff5502C6C096
// 账户地址：0xaECAC23B2bbB13748c7b507b0df8171061cE4955
// 私钥：955c8cb6be6243b8f2d84c0dcf81353c5e23fe3f83bf5878f012a55f9463e911

// 转发Matic币
// from 发送地址上转出
// to 接收的地址
// amountMatic Matic的数量
func (instance *PolygonClient) TransactionMatic(from *Wallet, to string, token string,amountMatic *big.Int) error {
	// 解析私钥
	privateKey, err := crypto.HexToECDSA(from.Key)
	if nil != err {
		return err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return err
		//log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	// 公钥转地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := instance.Client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}
	// 代币传输不需要传输ETH，因此将交易“值”设置为“0”。
	value := big.NewInt(0) // in wei (0 eth)
	gasPrice, err := instance.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	// 先将您要发送代币的地址存储在变量中。
	toAddress := common.HexToAddress(to)
	// 让我们将代币合约地址分配给变量。
	tokenAddress := common.HexToAddress(token)
	// 传递要调用合约中的函数名字
	transferFnSignature := []byte("transfer(address,uint256)")
	//生成函数签名的Keccak256哈希。 然后我们只使用前4个字节来获取方法ID。
	kState := crypto.NewKeccakState()
	kState.Write(transferFnSignature)
	methodID := kState.Sum(nil)[:4]
	fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb

	// 地址前补零，32字节
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	//fmt.Println(hexutil.Encode(paddedAddress)) // 0x0000000000000000000000004592d8f8d7b001e72cb26a73e4fa1806a51ac79d

	amount := new(big.Int).Mul(amountMatic, big.NewInt(params.Ether))
	//amount.SetString("100000000000000000", 10) // 1000 tokens
	// 数量前补零， 32字节
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	//fmt.Println(hexutil.Encode(paddedAmount)) // 0x00000000000000000000000000000000000000000000003635c9adc5dea00000
	//方法ID，地址，数量添加到 数据字段的字节片。
	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)
	// 计算gas费用上限
	gasLimit, err := instance.Client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &toAddress,
		Data: data,
	})
	if err != nil {
		return err
	}
	fmt.Println(gasLimit) // 23256
	// 构建交易事务类型
	inner := &types.LegacyTx{
		Nonce:    nonce,
		To:       &tokenAddress,
		Value:    value,
		Gas:      gasLimit * 3,	// gas费用不足，导致交易失败。 * 2处理。
		GasPrice: gasPrice,
		Data:     data,
	}
	tx := types.NewTx(inner)
	chainID, err := instance.Client.NetworkID(context.Background())
	if err != nil {
		return err
	}
	// 私钥对事物进行签名
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return err
	}
	// 广播交易
	err = instance.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex()) // tx sent: 0xa56316b637a94c4cc0331c73ef26389d6c097506d581073f927275e7a6ece0bc
	return nil
}