package cache

import "sync"

type shard struct {
	values    [][]byte       // 存储的数据
	idleIndex *queue[int]    // 保存存储数据中空闲节点的下标
	keysMap   map[uint32]int // key对应的slice index
	mu        sync.Mutex     // 锁
}

func newShard(capSize uint) *shard {
	return &shard{
		values:    make([][]byte, 0, capSize),
		idleIndex: newQueue[int](),
		keysMap:   make(map[uint32]int),
	}
}

func (s *shard) Get(key uint32) []byte {
	s.mu.Lock()
	defer s.mu.Unlock()
	if index, ok := s.keysMap[key]; ok {
		return s.values[index]
	}
	return nil
}

func (s *shard) Set(key uint32, value []byte) {
	s.mu.Lock()
	if index, ok := s.keysMap[key]; ok {
		s.values[index] = value
	} else {
		if s.idleIndex.Size() > 0 {
			index = s.idleIndex.Pop()
			s.values[index] = value
			s.keysMap[key] = index
		} else {
			s.values = append(s.values, value)
			s.keysMap[key] = len(s.values) - 1
		}
	}
	s.mu.Unlock()
}

func (s *shard) Remove(key uint32) {
	s.mu.Lock()
	if index, ok := s.keysMap[key]; ok {
		s.values[index] = nil
		s.idleIndex.Push(index)
		delete(s.keysMap, key)
	}
	s.mu.Unlock()
}
