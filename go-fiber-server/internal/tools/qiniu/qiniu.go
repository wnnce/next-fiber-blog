package qiniu

import "go-fiber-ent-web-layout/internal/conf"

var qiniuConfig conf.QiniuConfig

func init() {
	conf.CollectConfigReader(func(config *conf.Bootstrap) {
		qiniuConfig = config.Qiniu
	})
}
