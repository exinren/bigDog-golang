package global

import (
	"bigDog-golang/common"
	"bigDog-golang/pkg/logger"
	"bigDog-golang/pkg/setting"
)

// ćšć±ćé
var (
	ServerSetting		*setting.ServerSettingS
	AppSetting			*setting.AppSettingS
	Logger				*logger.Logger
	RedisSetting		*setting.RedisSettingS
	PolygonSetting		*setting.PolygonSettingS
	PolygonClients		*common.PolygonClient
	ERCContractClients	*common.ERCContractClient
)


