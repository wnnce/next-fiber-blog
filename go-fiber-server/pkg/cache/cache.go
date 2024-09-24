package cache

// 默认的分片大小
const defaultShardSize = 16

type memoryCache struct {
	shards    []*shard     // 分片
	shardSize uint32       // 分片数量
	hg        hashGenerate // key hash生成接口
}

func newMemoryCache(shardSize uint32) *memoryCache {
	var finalSize uint32
	if isPowerOfTwo(shardSize) {
		finalSize = shardSize
	} else {
		finalSize = defaultShardSize
	}
	shards := make([]*shard, 0, finalSize)
	for i := 0; i < int(finalSize); i++ {
		shards = append(shards, newShard(8))
	}
	return &memoryCache{
		shards:    shards,
		shardSize: finalSize,
		hg:        newFnvHashGenerate(),
	}
}

func (c *memoryCache) Set(key string, value []byte) error {
	hasher, err := c.hg.Generate([]byte(key))
	if err != nil {
		return err
	}
	index := hasher & (c.shardSize - 1)
	c.shards[index].Set(hasher, value)
	return nil
}

func (c *memoryCache) Get(key string) ([]byte, error) {
	hasher, err := c.hg.Generate([]byte(key))
	if err != nil {
		return nil, err
	}
	index := hasher & (c.shardSize - 1)
	return c.shards[index].Get(hasher), err
}

func (c *memoryCache) Remove(key string) {
	hasher, err := c.hg.Generate([]byte(key))
	if err != nil {
		return
	}
	index := hasher & (c.shardSize - 1)
	c.shards[index].Remove(hasher)
}

// 判断传入的数字是不是2的幂次方
func isPowerOfTwo(num uint32) bool {
	return (num & (num - 1)) == 0
}
