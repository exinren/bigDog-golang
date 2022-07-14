package tRedis

import (
	"bigDog-golang/pkg/setting"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

var redisPool *redis.Pool

func InitRedis(instance *setting.RedisSettingS){
	redisPool = &redis.Pool{
		MaxActive: 100,
		MaxIdle: 1,
		IdleTimeout: 3000,
		Wait: true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%v",instance.Address,instance.Port),
				redis.DialPassword(instance.Password),redis.DialDatabase(instance.Database))
		},
	}
}

//	SetString	存入string参数
//	key string	键
//	value value interface{} 值
//	expireSecond int 缓存时长，为0即为永久
//	bool	true成功，false失败
func SetSingle(key string,value interface{},expireSecond int) bool{
	conn := redisPool.Get()
	defer conn.Close()
	var err error
	if 0 != expireSecond {
		_,err = conn.Do("Set",key,value,"EX",expireSecond)
	}else{
		_,err = conn.Do("Set",key,value)
	}
	if nil != err {
		return false
	}
	return true
}


//	GetString
//	key string
//	string
func GetString(key string) string {
	conn := redisPool.Get()
	defer conn.Close()
	result,err := conn.Do("Get",key)
	if nil != err || nil == result {
		return ""
	}
	return string(result.([]byte))
}

//	DelKey	删除
//	key string
//	bool
func DelKey(key string) bool{
	conn := redisPool.Get()
	defer conn.Close()
	result,err := conn.Do("DEL",key)
	r := result.(int64)
	if nil != err || 0 == r{
		return false
	}
	return true
}

