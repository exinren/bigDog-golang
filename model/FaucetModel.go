package model

import (
	"bigDog-golang/global"
	"bigDog-golang/utils"
)

func SendCoin(address string) error {
	// 导入钱包
	wallet, err := global.PolygonClients.ImportWallet(global.PolygonSetting.PrivateKey)
	defer global.PolygonClients.CloseWallet(wallet)
	if nil != err {
		global.Logger.Error("导入钱包错误")
		return err
	}
	// 获取当前账号的余额
	count, err := global.ERCContractClients.GetBalanceOf(global.PolygonSetting.PublicKey)
	if nil != err {
		global.Logger.Error("查询当前账户余额代币错误")
		return err
	}
	// 计算发送的代币
	faucetValue := utils.CalcuFaucet(count)
	// 转发代币
	err = global.PolygonClients.TransactionMatic(wallet, address, global.PolygonSetting.ContractAddress,faucetValue)
	if nil != err {
		return err
	}
	return nil
}
