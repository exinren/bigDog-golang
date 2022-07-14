package routers

import (
	"bigDog-golang/router/ctrl"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	// 打印日志。
	r.Use(gin.Logger())
	// 处理panic
	r.Use(gin.Recovery())
	r.GET("/captcha", ctrl.Captcha)
	r.POST("/check", ctrl.CheckCaptchas)
	//r.GET("/captcha/verity/:value", func(c *gin.Context) {
	//	fmt.Println(123)
	//})

	//创建组
	//apiv1 := r.Group("api/v1")
	//apiv1.Use(middleware.JWT())
	// 标签相关的接口
	//{
	//	apiv1.POST("/tags", tag.Create)
	//	apiv1.DELETE("/tags/:id", tag.Delete)
	//	apiv1.PUT("/tags/:id", tag.Update)
	//	apiv1.PATCH("/tags/:id/state", tag.Update)
	//	apiv1.GET("/tags", tag.List)
	//}
	return r
}