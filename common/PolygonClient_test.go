package common

import (
	"bigDog-golang/pkg/setting"
	"math/big"
	"testing"
)

func TestConnNode(t *testing.T){
	//key := "73ae013380277a27e66842fbeee965b4b86fc66f2c555c7ecd350611a9619bb4"
	keyPath := "D://tmp"

	nodeURL := "https://rpc-mumbai.maticvigil.com/"

	//url := "https://evmtestnet.confluxrpc.com"
	polygon := setting.PolygonSettingS{
		NodeURL: nodeURL,
		KeyStorePath: keyPath,
	}
	client,err := InitPolygonClient(&polygon)
	//wallet,err := client.CreateWallet()
	//key := "955c8cb6be6243b8f2d84c0dcf81353c5e23fe3f83bf5878f012a55f9463e911"
	//t.Log(wallet)
	//client.Client
	//epoch, err := client.Client.GetEpochNumber()
	//t.Log(epoch)
	//client.CreateWallet("")
	t.Log("-----------------------------")
	//filePath := "./wallets/UTC--2022-07-12T07-52-47.254875000Z--11e6c161b236663b2146d35bcc3b28d592f936f9"
	wal,err := client.ImportWallet("955c8cb6be6243b8f2d84c0dcf81353c5e23fe3f83bf5878f012a55f9463e911")
	defer client.CloseWallet(wal)
	t.Log(wal)
	client.CloseWallet(wal)
	//t.Log(err)
	//t.Log(wal)
	t.Log(err)
	//t.Log(wal.Address.Address.String())
	//t.Log(123123123)

	// 合约地址：0xe9f4c1CB0F71B4c59612Dab5A0C0Ff5502C6C096
	// 账户地址：0xaECAC23B2bbB13748c7b507b0df8171061cE4955
	// 私钥：955c8cb6be6243b8f2d84c0dcf81353c5e23fe3f83bf5878f012a55f9463e911

	// 转发Matic币
	// from 发送地址上转出
	// to 接收的地址
	// amountMatic Matic的数量
	var (
		//key = "955c8cb6be6243b8f2d84c0dcf81353c5e23fe3f83bf5878f012a55f9463e911"
		to = "0x51a08D9797Ad0c58b205bff0AEB127eedF75b5D1"
		token = "0x587f61B649Bccd80d898017B518c15b7DE79b186"
		amountMatic = big.NewInt(10000)
	)
	err = client.TransactionMatic(wal,to,token,amountMatic)
	//t.Log(err)
}

