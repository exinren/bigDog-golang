package setting

import (
	"github.com/spf13/viper"
)

type Setting struct {
	vp *viper.Viper
}

func NewSetting () (*Setting, error) {
	vp := viper.New()
	//设定配置文件的名称为 config
	vp.SetConfigName("config-dev")
	//设置其配置路径为相对路径 configs/
	vp.AddConfigPath("./")
	//配置类型为 yaml
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if nil != err {
		return nil, err
	}
	return &Setting{vp}, nil
}