package common

import (
	"bigDog-golang/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
)
type ERCContractClient struct {
	ABIStr string
	Address *common.Address
	Client *ethclient.Client
	Contract *bind.BoundContract
}

/** @dev 初始化代币合约
*	abi： 合约的abi
*	addres： 合约部署的地址
*	client： 连接节点的客户端
 */
func InitERCContractClient(abi ,address string,client *ethclient.Client) (*ERCContractClient,error) {
	contractAddress := common.HexToAddress(address)
	contract, err := bindToken(abi, contractAddress,client)
	if nil != err {
		return nil,err
	}
	c := ERCContractClient{
		ABIStr: abi,
		Address: &contractAddress,
		Client: client,
		Contract: contract,
	}
	return &c,nil
}

// 绑定合约
func bindToken(contractABI string, address common.Address, backend bind.ContractBackend) (*bind.BoundContract, error) {
	// 解析默认abi
	parsed, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, backend, backend, backend), nil
}

// 通过地址查看余额
// address: 例如：0xaECAC23B2bbB13748c7b507b0df8171061cE4955
// big.Int， 返回的余额已扣除小数位
func (erc *ERCContractClient) GetBalanceOf(address string) (*big.Int, error) {
	owner := common.HexToAddress(address)
	var bal = make([]interface{}, 0)
	err := erc.Contract.Call(nil, &bal,"balanceOf", owner)
	if nil != err {
		return nil, err
	}
	// 转换类型
	balString := utils.TransType(bal[0])
	balCount := new(big.Int)
	big, ok := balCount.SetString(balString, 10)
	eth := utils.ToDecimal(big, 18)
	if !ok {
		return nil, err
	}
	return eth.BigInt(), nil
}
