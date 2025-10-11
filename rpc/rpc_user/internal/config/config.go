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
	Avatar struct {
		UploadPath string `json:",default=../../static/avatar"`
	}
	Email struct {
		Host     string
		Port     int
		User     string
		Password string
	}
	ProductRpcConf zrpc.RpcClientConf
}
