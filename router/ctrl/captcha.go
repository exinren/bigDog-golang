package ctrl

import (
	captchas "bigDog-golang/common/captcha"
	"bigDog-golang/constant"
	"bigDog-golang/model"
	"bigDog-golang/pkg/app"
	"bigDog-golang/pkg/errcode"
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
		res := gin.H{"code": 401,"msg": "传输参数错误"}
		responese.ToResponse(res)
		return
	}
	// 验证
	redis := tRedis.GetInstanceByCaptchaStore()
	value := redis.Get(fmt.Sprintf("captcha:%s", id),true)
	if value == "" {
		res := gin.H{"code": 401,"msg": "id不存在"}
		responese.ToResponse(res)
		return
	}
	// 转化数据
	valueInt, _ := strconv.Atoi(value)
	if (valueInt - leftInt) >= 10 || (leftInt - valueInt) >= 10{
		responese.ToErrorResponse(errcode.ErrorVerityCaptchaFail)
		return
	}
	if address == "" {
		res := gin.H{"code": 401,"msg": "地址不能为空"}
		responese.ToResponse(res)
		return
	}
	// 获取ip和address缓存并查询
	ip := c.ClientIP()
	ipKey := fmt.Sprintf("%s%s",constant.IP, ip)
	ipValue := redis.Get(ipKey, false)
	if ipValue == ip && ip != "" {
		res := gin.H{"code": 401,"msg": "领取失败，已经领取过了！"}
		responese.ToResponse(res)
		return
	}
	// 地址储存redis的key
	key := fmt.Sprintf("%s%s",constant.ADDRESS, address)
	keyValue := redis.Get(key, false)
	if keyValue == address {
		res := gin.H{"code": 401,"msg": "领取失败，已经领取过了！"}
		responese.ToResponse(res)
		return
	}
	err = model.SendCoin(address)
	if nil != err {
		res := gin.H{"code": 500,"msg": "领取失败，请稍后再试！"}
		responese.ToResponse(res)
		return
	}
	// 存储ip和address
	tRedis.SetSingle(key, address, constant.ONEDAY)
	tRedis.SetSingle(ipKey, ip, constant.ONEDAY)
	res := gin.H{"code": 200, "msg": "领取成功"}
	fmt.Println("ip:",ip)
	//fmt.Println("key:",key)
	responese.ToResponse(res)
	return
}