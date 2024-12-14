package logic

import (
	"context"
	"fmt"
	"gozerodtm/comm"
	"gozerodtm/stock-srv/internal/gen"
	"gozerodtm/stock-srv/internal/svc"
	"gozerodtm/stock-srv/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dtm-labs/client/dtmcli"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductLogic {
	return &DeductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeductLogic) Deduct(in *pb.DecuctReq) (*pb.DeductResp, error) {

	fmt.Printf("扣库存start....")
	uid := comm.GenId()
	OK, err := l.svcCtx.RedisCilent.SetnxExCtx(l.ctx, comm.Deduct, uid, 10)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !OK {
		return nil, status.Error(codes.Aborted, comm.NOKeyRedis)
	}

	stock, err := l.svcCtx.StockModel.FindOneByGoodsId(l.ctx, in.GoodsId)
	if err != nil && err != gen.ErrNotFound {
		err1 := l.DelRedisCtx(l.ctx, uid)
		logx.Error(fmt.Errorf(" err != nil && err != gen.ErrNotFound 解锁情况: %v Deduct err: %v ", err1, err))
		return nil, status.Error(codes.Internal, "解锁情况"+err.Error())
	}
	if stock == nil || stock.Num < in.Num {
		err1 := l.DelRedisCtx(l.ctx, uid)
		logx.Error(fmt.Errorf(" stock == nil || stock.Num < in.Num 解锁情况: %v Deduct err: %v ", err1, err))
		return nil, status.Error(codes.Aborted, comm.NumOutOfRangeCode)
	}

	sqlResult, err := l.svcCtx.StockModel.DecuctStock(l.ctx, in.GoodsId, in.Num)
	if err != nil {
		err1 := l.DelRedisCtx(l.ctx, uid)
		logx.Error(fmt.Errorf("DecuctStock(l.ctx, in.GoodsId, in.Num)  解锁情况: %v Deduct err: %v ", err1, err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	affected, err := sqlResult.RowsAffected()
	if err != nil {
		err1 := l.DelRedisCtx(l.ctx, uid)
		logx.Error(fmt.Errorf("sqlResult.RowsAffected() 解锁情况: %v Deduct err: %v ", err1, err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//如果是影响行数为0，直接就告诉dtm失败不需要重试了
	if affected <= 0 {
		err1 := l.DelRedisCtx(l.ctx, uid)
		logx.Error(fmt.Errorf("affected <= 0 解锁情况: %v Deduct err: %v ", err1, err))
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	//！！开启测试！！ ： 测试订单回滚更改状态为失效，并且当前库扣失败不需要回滚
	//return fmt.Errorf("扣库存失败 err : %v , in:%+v \n",err,in)
	err1 := l.DelRedisCtx(l.ctx, uid)
	logx.Error(fmt.Errorf("解锁情况: %v Deduct err: %v ", err1, err))
	return &pb.DeductResp{}, err1
}
func (l *DeductLogic) DelRedisCtx(ctx context.Context, uid string) error {
	_, err := l.svcCtx.RedisCilent.DelCtx(l.ctx, comm.Deduct, uid)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	return nil
}
