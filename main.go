package main

import (
	"bigDog-golang/common"
	"bigDog-golang/global"
	"bigDog-golang/pkg/logger"
	"bigDog-golang/pkg/setting"
	"bigDog-golang/pkg/timer"
	routers "bigDog-golang/router"
	"bigDog-golang/tRedis"
	"context"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 初始化
func init()  {
	// 初始化配置文件
	err := setupSetting()
	if nil != err {
		log.Fatalf("init.setupSetting err：%v", err)
	}
	// 初始化日志
	err = setupLogger()
	if nil != err {
		log.Fatalf("init.setupLogger err： %v", err)
	}
	// 初始化redis
	tRedis.InitRedis(global.RedisSetting)
	// 初始化扑该仔链
	initPolygon()
	// 初始化定时任务
	go timer.InitCron()
}

// @title bigDogCoin项目
// @version 1.0
// @description Go 做项目
func main()  {
	gin.SetMode(global.ServerSetting.RunMode)
	route := routers.NewRouter()
	
	serve := &http.Server{
		Addr: global.ServerSetting.HttpPort,
		Handler: route,
		ReadTimeout: global.ServerSetting.ReadTimeout,
		WriteTimeout: global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	// 优雅的关机
	//通过子程序启动服务
	go func() {
		log.Printf("Actual pid is %d", syscall.Getpid())
		if err := serve.ListenAndServe(); err != nil {
			//global.Logger.Info("Listen： %s\n", err)
			log.Printf("Listen： %s\n", err)
		}

	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	// 阻塞，当进程中断的时候唤醒，如： ctrl + c
	<- quit
	global.Logger.Info("Shutdown Server....")
	// 延迟十秒后重启
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	// 实现优雅的重启
	if err := serve.Shutdown(ctx); nil != err {
		log.Fatal("Server Shutdown：", err)
	}
	global.Logger.Info("Server exit")
}

// 读取配置文件
func setupSetting() error {
	s, err := setting.NewSetting()
	if nil != err {
		return err
	}
	// 读取Server配置
	err = s.ReadSection("Server", &global.ServerSetting)
	if nil != err {
		return err
	}
	// 读取App配置
	err = s.ReadSection("App", &global.AppSetting)
	if nil != err {
		return err
	}
	err = s.ReadSection("Redis",&global.RedisSetting)
	if nil != err {
		return err
	}
	err = s.ReadSection("Polygon", &global.PolygonSetting)
	if nil != err {
		return err
	}
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	return nil
}
// 日志初始化设置
func setupLogger () error {
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename: global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
		MaxSize: 600,
		MaxAge: 10,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)
	return nil
}

// 初始化代币
func initPolygon() error {
	var err error
	global.PolygonClients, err = common.InitPolygonClient(global.PolygonSetting)
	if nil != err {
		global.Logger.Error("初始化扑该仔错误")
		return err
	}
	global.ERCContractClients, err = common.InitERCContractClient(global.PolygonSetting.ContractABI, global.PolygonSetting.ContractAddress, global.PolygonClients.Client)
	if nil != err {
		global.Logger.Error("初始化大狗币合约错误")
		return err
	}
	return nil
}