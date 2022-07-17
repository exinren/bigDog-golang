package timer

import (
	"bigDog-golang/constant"
	"bigDog-golang/global"
	"bigDog-golang/tRedis"
	"github.com/robfig/cron"
	"time"
)

// 定时任务
func InitCron()  {
	c := cron.New()
	// 设置定时任务。可以添加多个
	c.AddFunc("0 0 0 * * ?", func() {
		clearCacheAddress()
	})
	c.Start()
	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)
		}
	}
}
// 清除领取过的账号和ip缓存
func clearCacheAddress()  {
	isSuccess := tRedis.ExpMatchDel(constant.ADDRESS)
	if isSuccess {
		global.Logger.Info("领取账号缓存清除成功!")
	} else {
		global.Logger.Info("领取账号缓存清除失败!")
	}
	isSuccess = tRedis.ExpMatchDel(constant.IP)
	if isSuccess {
		global.Logger.Info("领取IP缓存清除成功!")
	} else {
		global.Logger.Info("领取IP缓存清除失败!")
	}
}