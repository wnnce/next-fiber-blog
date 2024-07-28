package optimize

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/chai2010/webp"
	"hash"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strings"
	"sync"
)

var hasherPool *sync.Pool
var imageOnce sync.Once

func initPool() {
	hasherPool = &sync.Pool{
		New: func() any {
			return md5.New()
		},
	}
}

// CompressImageWithWebp 图片格式化为Webp格式
func CompressImageWithWebp(reader io.Reader, lossless bool, quality float32) (*bytes.Buffer, error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	buff := &bytes.Buffer{}
	err = webp.Encode(buff, img, &webp.Options{
		Lossless: lossless,
		Quality:  quality,
	})
	return buff, err
}

// CheckWebpFormatSupport 检查图片是否可以被转换为webp格式
// fileName 待转换的图片名称
func CheckWebpFormatSupport(fileName string) bool {
	supports := [3]string{".png", ".jpg", ".jpeg"}
	for i := 0; i < 3; i++ {
		if strings.HasSuffix(fileName, supports[i]) {
			return true
		}
	}
	return false
}

// ComputeMd5 计算Md5
func ComputeMd5(bytes []byte) string {
	imageOnce.Do(initPool)
	hasher := hasherPool.Get().(hash.Hash)
	defer func() {
		hasher.Reset()
		hasherPool.Put(hasher)
	}()
	hasher.Write(bytes)
	return hex.EncodeToString(hasher.Sum(nil))
}
