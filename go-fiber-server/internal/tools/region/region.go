package region

import (
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"go-fiber-ent-web-layout/internal/conf"
	"log/slog"
	"sync"
)

var (
	once     sync.Once
	searcher *xdb.Searcher
	filePath = "./ip2region.xdb"
)

func init() {
	conf.CollectConfigReader(readXdbConfig)
}

func readXdbConfig(config *conf.Bootstrap) {
	filePath = config.XdbFilePath
}

func initSearcher() {
	buff, err := xdb.LoadContentFromFile(filePath)
	if err != nil {
		slog.Error("读取Xdb文件失败", "err", err)
		panic(err)
	}
	// 使用Memory模式 确保并发安全
	searcher, err = xdb.NewWithBuffer(buff)
	if err != nil {
		slog.Error("初始化searcher失败", "err", err)
		panic(err)
	}
}

func SearchLocation(ip string) string {
	once.Do(initSearcher)
	search, err := searcher.SearchByStr(ip)
	if err != nil {
		slog.Error("查询IP属地失败", "err", err)
		return ""
	}
	return search
}
