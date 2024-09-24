package optimize

import (
	"fmt"
	"testing"
)

func TestCompressTextWithGzip(t *testing.T) {
	value, _ := CompressTextWithGzip("hello world")
	origin, _ := DecompressTextWithGzip(value)
	fmt.Println(value, origin)
}
