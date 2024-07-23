package cache

import (
	"fmt"
	"testing"
)

var qe = newQueue[int]()

func TestQueue_Push(t *testing.T) {
	for i := 0; i < 1000; i++ {
		qe.Push(i)
	}
	fmt.Println(qe.Size())
	for i := 0; i < 1000; i++ {
		fmt.Println(qe.Pop())
	}
	fmt.Println(qe.Size())
}

func BenchmarkQueue_Push(b *testing.B) {
	for i := 0; i < 1000; i++ {
		qe.Push(i)
	}
	fmt.Println(qe.Size())
	for i := 0; i < 1000; i++ {
		qe.Pop()
	}
	fmt.Println(qe.Size())
}
