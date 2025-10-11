package svc

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sk_mall/rpc/rpc_product/product"
	"sk_mall/rpc/rpc_user/internal/config"
	"sk_mall/utils"
	"time"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config          config.Config
	DBConn          sqlx.SqlConn
	Rds             *redis.Redis
	EmailSender     *utils.EmailSender
	EmailCodePrefix string
	ProductRpc      product.Product
}

func CreateAvatarUploadDir(path string) {
	//创建upload目录
	if !utils.IsDirExists(path) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			errStr := fmt.Sprintf("avatar upload 目录创建失败 err:%s", err.Error())
			logx.Severe(errStr)
			panic(errStr)
		}
	}
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := sql.Open("mysql", c.DB.DataSource)
	if err != nil {
		logc.Error(context.Background(), "数据库连接失败")
		panic(err)
	}

	if c.DB.MaxOpen > 0 {
		db.SetMaxOpenConns(c.DB.MaxOpen)
	}
	if c.DB.MaxIdle > 0 {
		db.SetMaxIdleConns(c.DB.MaxIdle)
	}
	if c.DB.MaxLifetime > 0 {
		db.SetConnMaxLifetime(time.Duration(c.DB.MaxLifetime) * time.Second)
	}

	// 包装成 sqlx.SqlConn
	conn := sqlx.NewSqlConnFromDB(db)

	//redis:
	rds := redis.MustNewRedis(c.Redis.RedisConf)
	redis_flag := rds.Ping()
	if !redis_flag {
		logc.Error(context.Background(), "[Redis]redis连接失败")
		panic("[Redis]redis连接失败")
	}

	//email
	var e utils.EmailSender
	e.InitGomali(c.Email.Host, c.Email.Port, c.Email.User, c.Email.Password)

	CreateAvatarUploadDir(c.Avatar.UploadPath)
	return &ServiceContext{
		Config:          c,
		DBConn:          conn,
		Rds:             rds,
		EmailSender:     &e,
		EmailCodePrefix: "Rpc_User:EmailCode:",
		ProductRpc: product.NewProduct(zrpc.MustNewClient(c.ProductRpcConf)),
	}
}
