package conf

import (
	"go-fiber-ent-web-layout/internal/middleware/limiter"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
	"time"
)

type ConfigReader func(config *Bootstrap)

var (
	depends []ConfigReader
	mu      sync.Mutex
)

type Bootstrap struct {
	Server      Server      `json:"server" yaml:"server"`
	Data        Data        `json:"data" yaml:"data"`
	Jwt         Jwt         `json:"jwt" yaml:"jwt"`
	Qiniu       QiniuConfig `json:"qiniu" yaml:"qiniu"`
	XdbFilePath string      `json:"xdbFilePath" yaml:"xdb-file-path"`
}

type Server struct {
	Name    string        `json:"name" yaml:"name"`
	Host    string        `json:"host" yaml:"host"`
	Port    uint          `json:"port" yaml:"port"`
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
	Limiter struct {
		Sliding     limiter.SlidingConfig     `json:"sliding" yaml:"sliding"`
		TokenBucket limiter.TokenBucketConfig `json:"tokenBuket" yaml:"token-buket"`
	} `json:"limiter" yaml:"limiter"` // 限流配置
}

type Data struct {
	Database struct {
		Driver   string `json:"driver" yaml:"driver"`
		Host     string `json:"host" yaml:"host"`
		Port     int    `json:"port" yaml:"port"`
		Username string `json:"username" yaml:"username"`
		Password string `json:"password" yaml:"password"`
		DbName   string `json:"dbName" yaml:"db-name"`
	} `json:"database" yaml:"database"`
	Redis struct {
		Host        string        `json:"host" yaml:"host"`
		Port        int           `json:"port" yaml:"port"`
		Index       int           `json:"index" yaml:"index"`
		Username    string        `json:"username" yaml:"username"`
		Password    string        `json:"password" yaml:"password"`
		ReadTimeout time.Duration `json:"readTimeout" yaml:"read-timeout"`
		WireTimeout time.Duration `json:"wireTimeout" yaml:"wire-timeout"`
	} `json:"redis" yaml:"redis"`
}

type Jwt struct {
	Issue      string        `json:"issue" yaml:"issue"`
	ExpireTime time.Duration `json:"expireTime" yaml:"expire-timeout"`
	Secret     string        `json:"secret" yaml:"secret"`
}

type QiniuConfig struct {
	AccessKey    string `json:"accessKey" yaml:"access-key"`
	SecretKey    string `json:"secretKey" yaml:"secret-key"`
	Bucket       string `json:"bucket" yaml:"bucket"`              // 存储桶名称
	Region       string `json:"region" yaml:"region"`              // 存储区域域名
	BucketDomain string `json:"bucketDomain" yaml:"bucket-domain"` // 存储桶加速域名
}

// ReadConfig 读取配置文件
// path 配置文件路径
func ReadConfig(path string) *Bootstrap {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	file, err := os.ReadFile(dir + path)
	if err != nil {
		panic(err)
	}
	bootstrap := &Bootstrap{}
	err = yaml.Unmarshal(file, bootstrap)
	if err != nil {
		panic(err)
	}
	return bootstrap
}

// CollectConfigReader 收集所有需要读取配置的函数
func CollectConfigReader(reader ConfigReader) {
	mu.Lock()
	depends = append(depends, reader)
	mu.Unlock()
}

// IssuedConfig 下发配置
func IssuedConfig(config *Bootstrap) {
	mu.Lock()
	for _, fn := range depends {
		fn(config)
	}
	mu.Unlock()
}
