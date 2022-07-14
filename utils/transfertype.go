package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
)

const MaxCount = 100000000 // 一亿
const MinCount = 10000000	// 一千万

var letterBytes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// 获取随机数
func RandStringBytesMaskImpr(n int) string {
	max := len(letterBytes)
	b := make([]rune, n)
	for i := range b {
		result,_ := rand.Int(rand.Reader, big.NewInt(int64(max)))
		b[i] = letterBytes[result.Int64()]
	}
	return string(b)
}

// interface转成对应类型的字符串。
func TransType(value interface{}) string {
	var key string
	if value == nil {
		return key
	}
	switch value.(type) {
	case int64:
		it := value.(int64)
		key = string(it)
	case int32:
		it := value.(int32)
		key = string(it)
	case string:
		key = value.(string)
	case bool:
		b := value.(bool)
		key = fmt.Sprintf("%v",b)
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}

/**
 * @Description: 获取根目录
 * @return string
 */
func GetPWD() string {
	path, err := os.Getwd()
	if err != nil {
		return ""
	}
	return path
}

func CalcuFaucet(total *big.Int) *big.Int {
	max := total.Int64()
	res := new(big.Int)
	if total.Int64() > MaxCount {
		max = MaxCount
	}
	result,_ := rand.Int(rand.Reader, big.NewInt(max))
	if (result.Int64() < MinCount) {
		return res.Add(result, big.NewInt(MinCount))
	}
	return result
}