package setting

import "time"

type ServerSettingS struct {
	RunMode string
	HttpPort string
	ReadTimeout time.Duration
	WriteTimeout time.Duration
}

type AppSettingS struct {
	DefaultPageSize int
	MaxPageSize int
	LogSavePath string
	LogFileName string
	LogFileExt string
}

type RedisSettingS struct {
	Address string
	Port int
	Password string
	Database int
}

type JWTSettingS struct {
	Secret	string
	Issuer	string
	Expire	time.Duration
}


func (s *Setting) ReadSection (k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if nil != err {
		return err
	}
	return nil
}