package comm

import "errors"

var NumOutOfRangeCode = errors.New("num out of range").Error()
var NumOutOfRange = errors.New("rpc error: code = Aborted desc = num out of range").Error()
var RedisEXPIRE = errors.New("lock expire").Error()

var Deduct = "DeducntStock"     // 扣减库存
var CreateOrder = "CreateOrder" // 创建订单
var NOKeyRedis = "no redisLock"
var RedisErr = "redis error"
