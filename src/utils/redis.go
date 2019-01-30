package utils

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"time"
)

type RedisClient struct {
	pool      *redis.Pool
	redisIp   string
	redisPort string
	redisPwd  string
	redisDb   int
}

func NewRedisClient(ip, port, pwd string, redis_db int) *RedisClient {
	connStr := fmt.Sprintf("%s:%s", ip, port)
	p := &redis.Pool{
		MaxIdle:     3,
		MaxActive:   10000,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", connStr)
			if err != nil {
				return nil, err
			}
			// 选择db
			if pwd != "" {
				logs.Info("redis use password:%s", pwd)
				c.Do("auth", pwd)
			} else {
				logs.Info("redis not use password")
			}
			c.Do("SELECT", redis_db)
			return c, nil
		},
	}

	c := &RedisClient{
		pool:      p,
		redisIp:   ip,
		redisPort: port,
		redisPwd: pwd,
		redisDb:   redis_db,
	}
	return c
}

func (this *RedisClient) GetConn() redis.Conn {
	rc := this.pool.Get()
	return rc
}
