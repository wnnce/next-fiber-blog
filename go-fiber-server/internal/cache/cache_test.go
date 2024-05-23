package cache

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"strconv"
	"sync"
	"testing"
	"time"
)

type student struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestSetAndGet(t *testing.T) {
	if err := Set("1", "11111", 1*time.Minute); err != nil {
		log.Fatalln(err)
	}
	if err := Set("1", "111111", 1*time.Minute); err != nil {
		log.Fatalln(err)
	}
}

func TestHash32(t *testing.T) {
	f := fnv.New32()
	fmt.Println(f.Size())
	_, _ = f.Write([]byte("hello"))
	fmt.Println(f.Size())
	fmt.Println(f.Sum32())
	fmt.Println(f.Size())
	f.Reset()
	_, _ = f.Write([]byte("world"))
	fmt.Println(f.Size())
	fmt.Println(f.Sum32())
	fmt.Println(f.Size())
}

func TestIsPowerOfTwo(t *testing.T) {
	fmt.Println(isPowerOfTwo(32))
	fmt.Println(isPowerOfTwo(4))
	fmt.Println(isPowerOfTwo(15))
	fmt.Println(isPowerOfTwo(24))
	fmt.Println(isPowerOfTwo(256))
}

func TestEncode(t *testing.T) {
	bytes, err := json.Marshal("hello")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(bytes)
	var str string
	if err = json.Unmarshal(bytes, &str); err != nil {
		log.Fatalln(err)
	}
	log.Println(str)
}

func TestFnvHashGenerate_Generate(t *testing.T) {
	hg := newFnvHashGenerate()
	for i := 0; i < 100; i++ {
		var wg sync.WaitGroup
		for y := 0; y < 10; y++ {
			go func() {
				wg.Add(1)
				defer wg.Done()
				key, err := hg.Generate([]byte(strconv.Itoa(i)))
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(key)
			}()
		}
		wg.Wait()
	}
}

func BenchmarkFnvHashGenerate_Generate(b *testing.B) {
	hg := newFnvHashGenerate()
	for i := 0; i < 100; i++ {
		var wg sync.WaitGroup
		for y := 0; y < 10; y++ {
			go func() {
				wg.Add(1)
				defer wg.Done()
				key, err := hg.Generate([]byte(strconv.Itoa(i)))
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(key)
			}()
		}
		wg.Wait()
	}
}

func BenchmarkCache_Set(t *testing.B) {
	ce := newMemoryCache(8)
	for i := 0; i < 10000; i++ {
		str := strconv.Itoa(i)
		if err := ce.Set(str, []byte(str)); err != nil {
			log.Fatalln(err)
		}
	}
	for i := 0; i < 10000; i++ {
		str := strconv.Itoa(i)
		value, err := ce.Get(str)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(value))
	}
	for i := 0; i < 10000; i++ {
		str := strconv.Itoa(i)
		ce.Remove(str)
	}
}

func TestTimestampToBytes(t *testing.T) {
	expire := 10 * time.Second
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(time.Now().Add(expire).UnixMilli()))
	fmt.Println(bytes)
	fmt.Println(len(bytes))
	num := binary.BigEndian.Uint64(bytes)
	fmt.Println(num)
}
