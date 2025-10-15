package svc

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sk_mall/rpc/rpc_product/internal/config"
	"sk_mall/utils"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	gr "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	DBConn    sqlx.SqlConn
	Rds       *redis.Redis
	RLCreater *utils.RedLockCreater
	ESCli     *elasticsearch.Client
}

func CreateCoverUploadDir(path string) {
	//创建upload目录
	if !utils.IsDirExists(path) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			errStr := fmt.Sprintf("Cover upload 目录创建失败 err:%s", err.Error())
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

	CreateCoverUploadDir(c.Cover.UploadPath)

	// 包装成 sqlx.SqlConn
	conn := sqlx.NewSqlConnFromDB(db)

	//redis:
	rds := redis.MustNewRedis(c.Redis.RedisConf)
	redis_flag := rds.Ping()
	if !redis_flag {
		logc.Error(context.Background(), "[Redis]redis连接失败")
		panic("[Redis]redis连接失败")
	}

	//RedLock
	var clis []*gr.Client
	for i, v := range c.DLockRedis.Hosts {
		r := gr.NewClient(&gr.Options{
			Addr:     v,
			Password: c.DLockRedis.Passes[i],
		})
		t := r.Ping(context.Background())
		if t.Err() != nil {
			logc.Errorf(context.Background(), "[RedLock]初始化失败,%s", t.Err().Error())
			panic("[RedLock]初始化失败")
		}
		clis = append(clis, r)
	}
	rlc, e3 := utils.NewRedLockCreater(clis)
	if e3 != nil {
		logc.Errorf(context.Background(), "[RedLock]初始化失败,%s", e3.Error())
		panic("[RedLock]初始化失败")
	}

	//Es
	cfg := elasticsearch.Config{
		Addresses: c.ES.Addr,
	}

	es, e4 := elasticsearch.NewClient(cfg)
	if e4 != nil {
		logc.Errorf(context.Background(), "[ES]初始化失败,%s", e4.Error())
		panic("[ES]初始化失败")
	}
	return &ServiceContext{
		Config:    c,
		DBConn:    conn,
		Rds:       rds,
		RLCreater: rlc,
		ESCli:     es,
	}
}
