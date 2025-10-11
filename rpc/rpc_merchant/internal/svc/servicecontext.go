package svc

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sk_mall/rpc/rpc_merchant/internal/config"
	"sk_mall/utils"
	"time"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	DBConn sqlx.SqlConn
}

func CreateLogoUploadDir(path string) {
	//创建upload目录
	if !utils.IsDirExists(path) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			errStr := fmt.Sprintf("Logo upload 目录创建失败 err:%s", err.Error())
			logx.Severe(errStr)
			panic(errStr)
		}
	}
}
func CreateLicenseUploadDir(path string) {
	//创建upload目录
	if !utils.IsDirExists(path) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			errStr := fmt.Sprintf("License upload 目录创建失败 err:%s", err.Error())
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

	CreateLogoUploadDir(c.Logo.UploadPath)
	CreateLicenseUploadDir(c.License.UploadPath)
	// 包装成 sqlx.SqlConn
	conn := sqlx.NewSqlConnFromDB(db)
	return &ServiceContext{
		Config: c,
		DBConn: conn,
	}
}
