package redis

import "github.com/zeromicro/go-zero/core/stores/redis"

var Redis *redis.Redis

func InitRedis(conf redis.RedisConf) {
	Redis = redis.MustNewRedis(conf)
}

func GetRedis() *redis.Redis {
	return Redis
}
