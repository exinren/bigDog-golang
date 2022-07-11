package ctrl

import (
	"bigDog-golang/pkg/app"
	"bigDog-golang/pkg/errcode"
	"github.com/gin-gonic/gin"
)

func  Captcha(c *gin.Context)  {
	responese := app.NewResponse(c)
	res := gin.H{"code": errcode.Success.Code(), "msg": errcode.Success.Msg()}
	responese.ToResponse(res)
	return
}