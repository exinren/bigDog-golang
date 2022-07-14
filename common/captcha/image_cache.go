package captchas

import (
	"bytes"
	"image"
	"io/ioutil"
	"path"
	"strings"
)

var (
	bgImgCache [][]byte //缓存背景图片
	bkImgCache [][]byte //缓存滑块模板图片
)

func LoadBackgroudImages(path string) (err error){
	bgImgCache, err = loadImages(path)
	return
}

func LoadBlockImages(path string) (err error) {
	bkImgCache, err = loadImages(path)
	return
}

func loadImages(basePath string) ([][]byte, error) {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return nil, err
	}
	var fileArr [][]byte
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasSuffix(f.Name(), ".png") {
			buf, err := ioutil.ReadFile(path.Join(basePath, f.Name()))
			if err != nil {
				return nil, err
			}
			fileArr = append(fileArr, buf)
		}
	}
	return fileArr, nil
}


// randBackgroudImage 随机抽取 背景图
func randBackgroudImage() (*ImageBuf, error) {
	n := r.Intn(len(bgImgCache))
	im, _, err := image.Decode(bytes.NewReader(bgImgCache[n]))
	if err != nil {
		return nil, err
	}
	return &ImageBuf{
		w: im.Bounds().Dx(),
		h: im.Bounds().Dy(),
		i: im,
	}, nil
}

// randBlockImage 抽取默认底图
func randBlockImage() (a *ImageBuf, err error) {
	n := r.Intn(len(bkImgCache))
	im, _, err := image.Decode(bytes.NewReader(bkImgCache[n]))
	if err != nil {
		return nil, err
	}
	a = &ImageBuf{
		w: im.Bounds().Dx(),
		h: im.Bounds().Dy(),
		i: im,
	}
	return
}
