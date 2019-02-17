package db

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

func init() {
	RedisPool = newPool(":6379")
}

var RedisPool *redis.Pool

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     1000,
		MaxActive:   1000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			// if _, err := c.Do("AUTH", password); err != nil {
			// 	c.Close()
			// 	return nil, err
			// }
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
