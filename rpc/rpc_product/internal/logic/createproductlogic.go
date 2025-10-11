package logic

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"sk_mall/rpc/rpc_product/internal/svc"
	"sk_mall/rpc/rpc_product/types"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductLogic {
	return &CreateProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateProductLogic) CreateProduct(in *__.CreateProductReq) (*__.CreateProductResp, error) {
	// todo: add your logic here and delete this line
	priceDec, err := decimal.NewFromString(in.Price)
	if err != nil {
		return &__.CreateProductResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "Price 解析错误",
			},
		}, nil
	}
	CoverName := uuid.NewString() + "." + in.CoverImageSuffix //api传后缀
	CoverPath := filepath.Join(l.svcCtx.Config.Cover.UploadPath, CoverName)

	//如果写文件失败头像为空字符串
	if e := os.WriteFile(CoverPath, in.CoverImage, 0644); e != nil {
		logc.Errorf(l.ctx, "[CoverUpload]文件写入失败: %s", e.Error())
		return &__.CreateProductResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "商品封面上传失败",
			},
		}, e
	}
	
	_, err = l.svcCtx.DBConn.ExecCtx(l.ctx, "insert into products(merchant_id,name,description,cover_image,price,stock,status) values(?,?,?,?,?,?,?)", in.MerchantId, in.Name, in.Desc, CoverPath, priceDec.String(), in.Stock, in.Status)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return &__.CreateProductResp{
					Base: &__.BaseResp{
						Code: 1,
						Msg:  fmt.Sprintf("商品 %s 已存在", in.Name),
					},
				}, nil
			} else {
				return &__.CreateProductResp{
					Base: &__.BaseResp{
						Code: 1,
						Msg:  "未知错误",
					},
				}, err
			}
		} else {
			return &__.CreateProductResp{
				Base: &__.BaseResp{
					Code: 1,
					Msg:  "未知错误",
				},
			}, err
		}
	}
	return &__.CreateProductResp{
		Base: &__.BaseResp{
			Code: 0,
			Msg:  "success",
		},
	}, nil
}
