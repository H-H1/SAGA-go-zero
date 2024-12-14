package svc

import (
	"gozerodtm/order-srv/internal/config"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	// "gozerodtm/order-srv/internal/model"
	"gozerodtm/order-srv/internal/gen"
)

type ServiceContext struct {
	Config      config.Config
	OrderModel  gen.OrderModel
	RedisCilent *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		OrderModel: gen.NewOrderModel(sqlx.NewMysql(c.DB.DataSource)),
		RedisCilent: redis.New(c.Cache.Host, func(r *redis.Redis) {
			r.Type = c.Cache.Type
			r.Pass = c.Cache.Pass
		}),
	}
}
