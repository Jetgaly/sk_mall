package svc

import (
	"context"
	"database/sql"
	"sk_mall/rpc/rpc_order/internal/config"
	"time"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	DBConn sqlx.SqlConn
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
	return &ServiceContext{
		Config: c,
		DBConn: conn,
	}
}
