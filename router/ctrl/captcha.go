package ctrl

import (
	captchas "bigDog-golang/common/captcha"
	"bigDog-golang/model"
	"bigDog-golang/pkg/app"
	"bigDog-golang/tRedis"
	"bigDog-golang/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Captcha(c *gin.Context)  {
	root := utils.GetPWD()
	bgPath := root + "/static/images/bg"
	blockPath := root + "/static/images/block"
	ret, err := captchas.Run(bgPath, blockPath)

	redis := tRedis.GetInstanceByCaptchaStore()
	if err != nil {
		return
	}
	randId := utils.RandStringBytesMaskImpr(16)
	id := fmt.Sprintf("captcha:%s", randId)
	redis.Set(id, fmt.Sprintf("%v",ret.Point.X))
	responese := app.NewResponse(c)
	res := gin.H{
		"id": randId,
		"y": ret.Point.Y,
		"im": ret.BackgroudImg,
		"imSlide": ret.BlockImg,
	}
	responese.ToResponse(res)
	return
}

func CheckCaptchas(c *gin.Context)  {
	id := c.PostForm("id")
	left := c.PostForm("left")
	address := c.PostForm("address")
	responese := app.NewResponse(c)
	leftInt, err := strconv.Atoi(left)

	if nil != err {
		res := gin.H{"code": 500,"info": "传输参数错误"}
		responese.ToResponse(res)
		return
	}
	redis := tRedis.GetInstanceByCaptchaStore()
	value := redis.Get(fmt.Sprintf("captcha:%s", id),true)
	if value == "" {
		res := gin.H{"code": 500,"info": "id不存在"}
		responese.ToResponse(res)
		return
	}
	valueInt, _ := strconv.Atoi(value)
	if (valueInt - leftInt) >= 10 || (leftInt - valueInt) >= 10{
		res := gin.H{"code": 504,"info": "验证失败"}
		responese.ToResponse(res)
		return
	}
	if address == "" {
		res := gin.H{"code": 500,"info": "address 不能为空"}
		responese.ToResponse(res)
		return
	}
	fmt.Println(address)
	err = model.SendCoin(address)
	if nil != err {
		res := gin.H{"code": 500,"info": "领取失败，请稍后再试！"}
		responese.ToResponse(res)
		return
	}
	res := gin.H{"code": 200, "info": "领取成功"}
	responese.ToResponse(res)
	return
}