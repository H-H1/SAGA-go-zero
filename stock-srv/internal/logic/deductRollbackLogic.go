package logic

import (
	"context"
	"fmt"

	"gozerodtm/stock-srv/internal/svc"
	"gozerodtm/stock-srv/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeductRollbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeductRollbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductRollbackLogic {
	return &DeductRollbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeductRollbackLogic) DeductRollback(in *pb.DecuctReq) (*pb.DeductResp, error) {

	fmt.Printf("库存回滚 in : %+v \n", in)

	if err := l.svcCtx.StockModel.AddStock(l.ctx, in.GoodsId, in.Num); err != nil {
		return nil, fmt.Errorf("回滚库存失败 err : %v ,goodsId:%d , num :%d", err, in.GoodsId, in.Num)
	}

	return &pb.DeductResp{}, nil
}
