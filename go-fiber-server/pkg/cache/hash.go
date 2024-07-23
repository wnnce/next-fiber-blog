package cache

import (
	"hash"
	"hash/fnv"
	"sync"
)

type hashGenerate interface {
	Generate(value []byte) (uint32, error)
}

func newFnvHashGenerate() hashGenerate {
	return &fnvHashGenerate{
		pool: &sync.Pool{
			New: func() any {
				return fnv.New32()
			},
		},
	}
}

type fnvHashGenerate struct {
	pool *sync.Pool
}

func (f *fnvHashGenerate) Generate(value []byte) (uint32, error) {
	hs := f.pool.Get().(hash.Hash32)
	defer func() {
		hs.Reset()
		f.pool.Put(hs)
	}()
	if _, err := hs.Write(value); err != nil {
		return 0, err
	}
	return hs.Sum32(), nil
}
