package logic

import (
	"context"
	"errors"
	"fmt"
	"gozerodtm/order-srv/order"
	"gozerodtm/stock-srv/stock"
	"net/http"
	"time"

	"gozerodtm/order-api/internal/svc"
	"gozerodtm/order-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateLogic {
	return CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req types.QuickCreateReq, r *http.Request) (*types.QuickCreateResp, error) {

	var err1 error
	var err2 error
	createOrderReq := &order.CreateReq{UserId: req.UserId, GoodsId: req.GoodsId, Num: req.Num}
	deductReq := &stock.DecuctReq{GoodsId: req.GoodsId, Num: req.Num}
	_, err1 = l.svcCtx.OrderRpc.Create(l.ctx, createOrderReq)
	if err1 != nil {
		return nil, err1
	}
	fmt.Println("err1---", err1)
	_, err2 = l.svcCtx.StockRpc.Deduct(l.ctx, deductReq)
	fmt.Println("err2---", err2)
	if err2 != nil {
		err := l.RetryRollback(l.svcCtx.OrderRpc.CreateRollback, createOrderReq, err1)
		return nil, errors.New(err2.Error() + " and " + err.Error())
	}
	var Resp *types.QuickCreateResp
	return Resp, nil
}

func (l *CreateLogic) RetryRollback(Rollback order.CreateRollbackfunc, createOrderReq *order.CreateReq, err1 error) error {
	if err1 == nil {
		retry := l.svcCtx.Config.Retry
		for retry > 0 {
			_, err := Rollback(l.ctx, createOrderReq)
			if err != nil {
				_, err = Rollback(l.ctx, createOrderReq)
				retry--
				if retry == 0 {
					logx.Error(fmt.Sprintf("deduct err and rollback failed: %v", err))
					return errors.New("deduct err and rollback failed")
				}
				continue
			}
			retry = 0
			logx.Info(fmt.Sprintf("deduct err and rollback successful: %v", err))
			return errors.New("deduct err and rollback successful")
		}
	}
	return nil
}
func Select(c context.Context, ch chan string, finish chan int) (resp *types.QuickCreateResp, err error) {
	select {
	case errmsg := <-ch:
		return &types.QuickCreateResp{}, errors.New(errmsg)
	case <-finish:
		return &types.QuickCreateResp{}, nil
	case <-time.After(5 * time.Second):

		return &types.QuickCreateResp{}, errors.New("timeout")
	}
}
