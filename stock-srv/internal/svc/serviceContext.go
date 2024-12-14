package svc

import (
	"gozerodtm/stock-srv/internal/config"
	"gozerodtm/stock-srv/internal/gen"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config      config.Config
	StockModel  gen.StockModel
	RedisCilent *redis.Redis // TODO: add redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		StockModel: gen.NewStockModel(sqlx.NewMysql(c.DB.DataSource)),
		RedisCilent: redis.New(c.Cache.Host, func(r *redis.Redis) {
			r.Type = c.Cache.Type
			r.Pass = c.Cache.Pass
		}),
	}
}
