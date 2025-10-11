package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	DB struct {
		DataSource  string
		MaxOpen     int
		MaxIdle     int
		MaxLifetime int
	}
	Cover struct {
		UploadPath string `json:",default=../../static/product_cover"`
	}
	DLockRedis struct {
		Hosts  []string
		Passes []string
	}
}
