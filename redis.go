package base

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

var RedisMaxActive, _ = strconv.Atoi(os.Getenv("RedisMaxActive"))
var _redis *RedisCache

type RedisCache struct {
	Default *redis.Pool
	Live    *redis.Pool
}

/**
 * redis_type [default默认 live直播间]
 */
func Redis() (*RedisCache, error) {
	var err error
	if _redis == nil {
		_redis = &RedisCache{}
		err = _redis.init(RedisConn)
	}
	return _redis, err
}

func (me *RedisCache) init(strType string) error {
	arr := strings.Split(strType, ",")
	for _, v := range arr {
		if err := me.open(strings.TrimSpace(v)); err != nil {
			return err
		}
	}
	return nil
}

func (me *RedisCache) open(redis_type string) error {
	switch redis_type {
	case "default":
		me.Default = me.getRedis(os.Getenv("REDIS_URL"))
	case "live":
		me.Live = me.getRedis(os.Getenv("LIVE_REDIS_URL"))
	default:
		return errors.New("redis type err")
	}

	return nil
}

func (me *RedisCache) Close() {
	if me.Default != nil {
		me.Default.Close()
	}

	if me.Live != nil {
		me.Live.Close()
	}
}

func (me *RedisCache) getRedis(url string) *redis.Pool {
	//log.WithFields(log.Fields{"url": url}).Error("getRedis")
	if RedisMaxActive == 0 {
		RedisMaxActive = 40
	}
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   RedisMaxActive,
		IdleTimeout: 240 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(url)
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
