package captchas

import (
	"bigDog-golang/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)
// 需要存储在redis
var SlideCaptchaData map[string]int

func NewSlideCaptcha(w http.ResponseWriter, wr *http.Request) {
	ret, err := Run("./a", "./b")
	if err != nil {
		return
	}
	id := utils.RandStringBytesMaskImpr(16)
	SlideCaptchaData[id] = ret.Point.X
	bt, _ := json.Marshal(map[string]interface{}{
		"id": id,
		"y": ret.Point.Y,
		"im": ret.BackgroudImg,
		"imSlide": ret.BlockImg,
	})

	if err != nil {
		return
	}
	w.Write(bt)
}

/**
	验证滑块的，带有凹块的
 */
func CheckSlideCaptcha(w http.ResponseWriter, r *http.Request)  {
	_ = r.ParseForm()
	id := r.Form.Get("id")
	left := r.Form.Get("left")
	leftInt, err := strconv.Atoi(left)
	if nil != err {
		bt, _ := json.Marshal(map[string]interface{}{
			"code": 500,
			"info": "传输参数错误",
		})
		_, _ = fmt.Fprintf(w, string(bt))
		return
	}
	if _, ok := SlideCaptchaData[id]; !ok {
		bt, _ := json.Marshal(map[string]interface{}{
			"code": 500,
			"info": "id不存在",
		})
		_, _ = fmt.Fprintf(w, string(bt))
		return
	}

	if (SlideCaptchaData[id] - leftInt) >= 10 {
		bt, _ := json.Marshal(map[string]interface{}{
			"code": 504,
			"info": "验证失败",
		})
		_, _ = fmt.Fprintf(w, string(bt))
		return
	}
	bt, _ := json.Marshal(map[string]interface{}{
		"code": 200,
		"info": "",
	})
	_, _ = fmt.Fprintf(w, string(bt))
	return
}
