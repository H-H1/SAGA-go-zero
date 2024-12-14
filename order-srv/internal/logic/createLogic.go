package logic

import (
	"context"

	"fmt"
	"gozerodtm/comm"
	"gozerodtm/order-srv/internal/gen"

	"gozerodtm/order-srv/internal/svc"
	"gozerodtm/order-srv/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateLogic) Create(in *pb.CreateReq) (*pb.CreateResp, error) {

	fmt.Printf("创建订单 in : %+v \n", in)
	uid := comm.GenId()
	OK, err := l.svcCtx.RedisCilent.SetnxExCtx(l.ctx, comm.CreateOrder, uid, 10)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !OK {
		return nil, status.Error(codes.Aborted, comm.NOKeyRedis)
	}
	order := new(gen.Order)
	order.GoodsId = in.GoodsId
	order.Num = in.Num
	order.UserId = in.UserId

	_, err = l.svcCtx.OrderModel.Insert(l.ctx, order)
	fmt.Printf("InsertInsertInsert创建订单 in : %+v \n", in)
	if err != nil {
		err1 := l.DelRedisCtx(l.ctx, uid)
		logx.Error(fmt.Errorf("解锁情况: %v 创建订单失败 err: %v order: %+v", err1, err, order))
		return nil, fmt.Errorf("解锁情况: %v 创建订单失败 err: %v order: %+v", err1, err, order)
	}
	err1 := l.DelRedisCtx(l.ctx, uid)
	fmt.Println("err1:", err1)
	logx.Error(fmt.Errorf("解锁情况: %v Deduct err: %v ", err1, err))
	if err1 != nil {
		logx.Error(fmt.Errorf("解锁情况: %v Deduct err: %v ", err1, err))
		return nil, err1
	}

	return &pb.CreateResp{}, err1
}

func (l *CreateLogic) DelRedisCtx(ctx context.Context, uid string) error {
	_, err := l.svcCtx.RedisCilent.DelCtx(l.ctx, comm.CreateOrder, uid)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	return nil
}
