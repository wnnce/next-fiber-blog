package optimize

import (
	"bytes"
	"compress/gzip"
	"sync"
)

var bufferPool *sync.Pool

func init() {
	bufferPool = &sync.Pool{
		New: func() any {
			return &bytes.Buffer{}
		},
	}
}

// CompressTextWithGzip 使用Gzip压缩文本
func CompressTextWithGzip(text string) (string, error) {
	buffer := bufferPool.Get().(*bytes.Buffer)
	defer func() {
		buffer.Reset()
		bufferPool.Put(buffer)
	}()
	writer := gzip.NewWriter(buffer)
	_, err := writer.Write([]byte(text))
	if err != nil {
		return "", err
	}
	err = writer.Close()
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// DecompressTextWithGzip 解压被gzip压缩的文本
func DecompressTextWithGzip(compressed string) (string, error) {
	buffer := bufferPool.Get().(*bytes.Buffer)
	buffer.WriteString(compressed)
	defer func() {
		buffer.Reset()
		bufferPool.Put(buffer)
	}()
	reader, err := gzip.NewReader(buffer)
	if err != nil {
		return "", err
	}
	defer reader.Close()
	out := bufferPool.Get().(*bytes.Buffer)
	defer func() {
		out.Reset()
		bufferPool.Put(out)
	}()
	_, err = out.ReadFrom(reader)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
