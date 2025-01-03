package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	//Consul consul.Conf //if consul use
	DB struct {
		DataSource string
	}
	Cache redis.RedisConf
}
