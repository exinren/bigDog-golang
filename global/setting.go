package global

import (
	"bigDog-golang/pkg/logger"
	"bigDog-golang/pkg/setting"
)

// 全局变量
var (
	ServerSetting		*setting.ServerSettingS
	AppSetting			*setting.AppSettingS
	Logger				*logger.Logger
	RedisSetting		*setting.RedisSettingS
	JWTSetting			*setting.JWTSettingS
)

