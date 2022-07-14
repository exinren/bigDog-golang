package captchas

import (
	"math/rand"
	"time"
)

// Point 随机生成的抠图位置
type Point struct {
	X int
	Y int
}

// CutoutRet 抠图出来的结果
type CutoutRet struct {
	Point        *Point
	BackgroudImg string
	BlockImg     string
}


var r = rand.New(rand.NewSource(time.Now().UnixNano()))
