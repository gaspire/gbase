package base

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

// OneMonthSeconds 缓存时长
const OneMonthSeconds = 30 * 24 * 60 * 60

var (
	// 缓存
	_redisCache *Redis
	// 微信
	_redisWechat      *Redis
	_redisDBCache, _  = strconv.Atoi(os.Getenv("REDIS_DB_CACHE"))
	_redisDBWechat, _ = strconv.Atoi(os.Getenv("REDIS_DB_WECHAT"))
)

// Redis redis 封装类
type Redis struct {
	db   int
	conn *redis.Pool
}

//RedisOpts redis 连接属性
type RedisOpts struct {
	Host        string `yml:"host" json:"host"`
	Password    string `yml:"password" json:"password"`
	Database    int    `yml:"database" json:"database"`
	MaxIdle     int    `yml:"max_idle" json:"max_idle"`
	MaxActive   int    `yml:"max_active" json:"max_active"`
	IdleTimeout int32  `yml:"idle_timeout" json:"idle_timeout"` //second
}

//NewRedisCache 实例化
func NewRedisCache() *Redis {
	if _redisCache == nil {
		if _redisDBCache == 0 {
			_redisDBCache = 2
		}
		_redisCache = NewRedis(_redisDBCache)
	}
	return _redisCache
}

//NewRedisWechat 实例化
func NewRedisWechat() *Redis {
	if _redisWechat == nil {
		if _redisDBWechat == 0 {
			_redisDBWechat = 4
		}
		_redisWechat = NewRedis(_redisDBWechat)
	}
	return _redisWechat
}

//NewRedis 实例化
func NewRedis(db int) *Redis {
	opts := &RedisOpts{
		Host:        os.Getenv("REDIS_HOST"),
		Password:    os.Getenv("REDIS_PASSWORD"),
		MaxIdle:     3,
		IdleTimeout: 240,
		Database:    db,
	}

	if opts.MaxActive == 0 {
		opts.MaxActive = 40
	}
	pool := &redis.Pool{
		MaxActive:   opts.MaxActive,
		MaxIdle:     opts.MaxIdle,
		IdleTimeout: time.Second * time.Duration(opts.IdleTimeout),
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", opts.Host,
				redis.DialDatabase(opts.Database),
				redis.DialPassword(opts.Password),
			)
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}
	return &Redis{conn: pool, db: opts.Database}
}

//Close 关闭
func (me *Redis) Close() {
	if me.conn != nil {
		me.conn.Close()
	}
}

//GetConn 设置conn
func (me *Redis) GetConn() *redis.Pool {
	return me.conn
}

//SetConn 设置conn
func (me *Redis) SetConn(conn *redis.Pool) {
	me.conn = conn
}

// SetDB 设置当前数据库
func (me *Redis) SetDB(db int) (val string) {
	conn := me.conn.Get()
	defer conn.Close()

	conn.Do("SELECT", db)
	return
}

// GetStr 获取字符串
func (me *Redis) GetStr(key string) (val string, err error) {
	conn := me.conn.Get()
	defer conn.Close()

	val, err = redis.String(conn.Do("GET", key))
	return
}

// SetStr 保存
func (me *Redis) SetStr(key string, val string) (err error) {
	conn := me.conn.Get()
	defer conn.Close()

	_, err = conn.Do("SET", key, val, "EX", OneMonthSeconds)
	return
}

// SetStrEx 保存
func (me *Redis) SetStrEx(key string, val string, expire int) (err error) {
	conn := me.conn.Get()
	defer conn.Close()

	_, err = conn.Do("SET", key, val, "EX", expire)
	return
}

// GetInt 获取整形
func (me *Redis) GetInt(key string) (val int) {
	conn := me.conn.Get()
	defer conn.Close()

	val, _ = redis.Int(conn.Do("GET", key))
	return
}

// Del 删除
func (me *Redis) Del(key string) (err error) {
	conn := me.conn.Get()
	defer conn.Close()

	_, err = conn.Do("DEL", key)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Del]")
	}
	return
}

//DelKeys 删除多个key
func (me *Redis) DelKeys(keys []string) (err error) {
	conn := me.conn.Get()
	defer conn.Close()

	values := redis.Args{}.AddFlat(keys)
	_, err = conn.Do("DEL", values...)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[DelKeys]")
	}
	return
}

// SetInt 保存整形
func (me *Redis) SetInt(key string, val int) (err error) {
	conn := me.conn.Get()
	defer conn.Close()

	_, err = conn.Do("SET", key, val)
	return
}

// IsExist 检查是否存在
func (me *Redis) IsExist(key string) (exist bool) {
	conn := me.conn.Get()
	defer conn.Close()

	exist, _ = redis.Bool(conn.Do("EXISTS", key))
	return
}

// SetEx 检查是否存在
func (me *Redis) SetEx(key string, times int) (err error) {
	conn := me.conn.Get()
	defer conn.Close()

	_, err = conn.Do("SETEX", key, times, 1)
	return
}

//HDel 删除hash值
func (me *Redis) HDel(key string, args ...interface{}) (err error) {
	conn := me.conn.Get()
	defer conn.Close()

	args = append([]interface{}{key}, args...)
	_, err = conn.Do("HDEL", args...)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[HDel]")
	}
	return
}

//HExists 判断是否存在
func (me *Redis) HExists(key string, field string) (exsist int) {
	conn := me.conn.Get()
	defer conn.Close()

	var err error
	exsist, err = redis.Int(conn.Do("HEXISTS", key, field))
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[HExists]")
	}
	return
}

//HIncrby 自增
func (me *Redis) HIncrby(key string, field string, increment int) (err error) {
	conn := me.conn.Get()
	defer conn.Close()

	_, err = conn.Do("HINCRBY", key, field, increment)
	return
}

func (me *Redis) hGet(key string, field string) (reply interface{}, err error) {
	conn := me.conn.Get()
	defer conn.Close()

	reply, err = conn.Do("HGET", key, field)
	if reply == nil {
		err = errors.New("not found")
		return
	}
	return
}

//HGetStr 获取hash值
func (me *Redis) HGetStr(key string, field string) (val string, err error) {
	var reply interface{}
	reply, err = me.hGet(key, field)
	if err != nil {
		return
	}

	val, err = redis.String(reply, err)
	return
}

// HGetInt 获取hash值
func (me *Redis) HGetInt(key string, field string) (val int, err error) {
	var reply interface{}
	reply, err = me.hGet(key, field)
	if err != nil {
		return
	}

	val, err = redis.Int(reply, err)
	return
}

//HMset 设置hash值
func (me *Redis) HMset(key string, args ...interface{}) (err error) {
	conn := me.conn.Get()
	defer conn.Close()

	args = append([]interface{}{key}, args...)
	_, err = conn.Do("HMSET", args...)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[HMset]")
	}
	return
}

//HGetAll 返回hash值
func (me *Redis) HGetAll(key string, val interface{}) (err error) {
	conn := me.conn.Get()
	defer conn.Close()

	v, err := redis.Values(conn.Do("HGETALL", key))
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[HGetAll]1")
		return
	}

	if err = redis.ScanStruct(v, val); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[HGetAll]2")
		return
	}
	return
}

//ZrangeInts 通过索引区间返回有序集合成指定区间内的成员
func (me *Redis) ZrangeInts(key string, start, stop int) (val []int) {
	conn := me.conn.Get()
	defer conn.Close()

	var err error
	val, err = redis.Ints(conn.Do("ZRANGE", key, start, stop))
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[ZrangeInts]")
	}
	return
}

//ZCount 通过索引区间返回有序集合成指定区间内的成员
func (me *Redis) ZCount(key string, min, max int) (val int) {
	conn := me.conn.Get()
	defer conn.Close()

	var err error
	val, err = redis.Int(conn.Do("ZCOUNT", key, min, max))
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[ZCOUNT]")
	}
	return
}

//ZrangeStrs 通过索引区间返回有序集合成指定区间内的成员
func (me *Redis) ZrangeStrs(key string, start, stop int) (val []string) {
	conn := me.conn.Get()
	defer conn.Close()

	var err error
	val, err = redis.Strings(conn.Do("ZRANGE", key, start, stop))
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[ZrangeStrs]")
	}
	return
}

//LrangeInts 获取列表指定范围内的元素
func (me *Redis) LrangeInts(key string, start, stop int) (val []int) {
	conn := me.conn.Get()
	defer conn.Close()

	var err error
	val, err = redis.Ints(conn.Do("LRANGE", key, start, stop))
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[LrangeInts]")
	}
	return
}

//LrangeStrs 获取列表指定范围内的元素
func (me *Redis) LrangeStrs(key string, start, stop int) (val []string) {
	conn := me.conn.Get()
	defer conn.Close()

	var err error
	val, err = redis.Strings(conn.Do("LRANGE", key, start, stop))
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[LrangeStrs]")
	}
	return
}

//LremInt 移除列表中与参数 VALUE 相等的元素
func (me *Redis) LremInt(key string, count, value int) (err error) {
	conn := me.conn.Get()
	defer conn.Close()

	_, err = conn.Do("LREM", key, count, value)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[LremInt]")
	}
	return
}

func (me *Redis) lPop(key string) (reply interface{}, err error) {
	conn := me.conn.Get()
	defer conn.Close()

	reply, err = conn.Do("LPOP", key)
	if reply == nil {
		err = errors.New("not found")
		return
	}
	return
}

//LPOPInt 移出并获取列表的第一个元素
func (me *Redis) LPOPInt(key string) (val int, err error) {
	var reply interface{}
	reply, err = me.lPop(key)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[LPOP]")
		return
	}
	val, err = redis.Int(reply, err)
	return
}

//RpopLpush right pop left push
func (me *Redis) RpopLpush(key string) (reply interface{}, err error) {
	conn := me.conn.Get()
	defer conn.Close()

	reply, err = conn.Do("RPOPLPUSH", key, key)
	if reply == nil {
		err = errors.New("not found")
		return
	}
	return
}

//RpopLpushInt 移除列表的最后一个元素，并将该元素添加到另一个列表并返回
func (me *Redis) RpopLpushInt(key string) (val int, err error) {
	var reply interface{}
	reply, err = me.RpopLpush(key)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[RPOPLPUSH]")
		return
	}
	val, err = redis.Int(reply, err)
	return
}

//RPush 向存于 key 的列表的尾部插入所有指定的值
func (me *Redis) RPush(key string, args []int) (err error) {
	conn := me.conn.Get()
	defer conn.Close()

	values := redis.Args{}.Add(key).AddFlat(args)
	_, err = conn.Do("RPUSH", values...)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[RPush]")
	}
	return
}
