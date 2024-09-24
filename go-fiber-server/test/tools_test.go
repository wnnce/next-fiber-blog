package test

import (
	"fmt"
	"testing"
)

func TestSprintf(t *testing.T) {
	fmt.Printf("'%%', %s", "hello")
}
