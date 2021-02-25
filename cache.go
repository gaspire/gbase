package gbase

import (
	"errors"

	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

// OneMonthSeconds 缓存时长
const OneMonthSeconds = 30 * 24 * 60 * 60

// Cache 缓存
type Cache struct {
	redisConn redis.Conn
}

var cache *Cache

//NewCache 初始化
func NewCache() (cache *Cache) {
	if cache == nil {
		cache = &Cache{}
		cache.init()
	}
	return cache
}

// init 初始化
func (me *Cache) init() (err error) {
	var redisObj *RedisCache
	redisObj, err = Redis()
	if err != nil {
		log.WithFields(log.Fields{"redisObj": redisObj, "err": err.Error()}).Error("init")
	}
	me.redisConn = redisObj.Default.Get()
	me.redisConn.Do("SELECT", 0)
	return
}

//Close 关闭
func (me *Cache) Close() {
	if me.redisConn == nil {
		return
	}
	me.redisConn.Close()
	//log.WithFields(log.Fields{"msg": "close"}).Info("[CACHE]")
}

// SetDB 设置当前数据库
func (me *Cache) SetDB(db int) (val string) {
	me.redisConn.Do("SELECT", db)
	return
}

// GetStr 获取字符串
func (me *Cache) GetStr(key string) (val string, err error) {
	val, err = redis.String(me.redisConn.Do("GET", key))
	return
}

// SetStr 保存
func (me *Cache) SetStr(key string, val string) (err error) {
	_, err = me.redisConn.Do("SET", key, val, "EX", OneMonthSeconds)
	return
}

// SetStr 保存
func (me *Cache) SetStrEx(key string, val string, expire int) (err error) {
	_, err = me.redisConn.Do("SET", key, val, "EX", expire)
	return
}

// GetInt 获取整形
func (me *Cache) GetInt(key string) (val int) {
	val, _ = redis.Int(me.redisConn.Do("GET", key))
	return
}

// Del 删除
func (me *Cache) Del(key string) (err error) {
	_, err = me.redisConn.Do("DEL", key)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[Del]")
	}
	return
}

//DelKeys 删除多个key
func (me *Cache) DelKeys(keys []string) (err error) {
	values := redis.Args{}.AddFlat(keys)
	_, err = me.redisConn.Do("DEL", values...)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[DelKeys]")
	}
	return
}

// SetInt 保存整形
func (me *Cache) SetInt(key string, val int) (err error) {
	_, err = me.redisConn.Do("SET", key, val)
	return
}

// IsExist 检查是否存在
func (me *Cache) IsExist(key string) (exist bool) {
	exist, _ = redis.Bool(me.redisConn.Do("EXISTS", key))
	return
}

// SetEx 检查是否存在
func (me *Cache) SetEx(key string, times int) (err error) {
	_, err = me.redisConn.Do("SETEX", key, times, 1)
	return
}

//HDel 删除hash值
func (me *Cache) HDel(key string, args ...interface{}) (err error) {
	args = append([]interface{}{key}, args...)
	_, err = me.redisConn.Do("HDEL", args...)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[HDel]")
	}
	return
}

//HExists 判断是否存在
func (me *Cache) HExists(key string, field string) (exsist int) {
	var err error
	exsist, err = redis.Int(me.redisConn.Do("HEXISTS", key, field))
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[HExists]")
	}
	return
}

//HIncrby 自增
func (me *Cache) HIncrby(key string, field string, increment int) (err error) {
	_, err = me.redisConn.Do("HINCRBY", key, field, increment)
	return
}

func (me *Cache) hGet(key string, field string) (reply interface{}, err error) {
	reply, err = me.redisConn.Do("HGET", key, field)
	if reply == nil {
		err = errors.New("not found")
		return
	}
	return
}

//HGetStr 获取hash值
func (me *Cache) HGetStr(key string, field string) (val string, err error) {
	var reply interface{}
	reply, err = me.hGet(key, field)
	if err != nil {
		return
	}

	val, err = redis.String(reply, err)
	return
}

// HGetInt 获取hash值
func (me *Cache) HGetInt(key string, field string) (val int, err error) {
	var reply interface{}
	reply, err = me.hGet(key, field)
	if err != nil {
		return
	}

	val, err = redis.Int(reply, err)
	return
}

//HMset 设置hash值
func (me *Cache) HMset(key string, args ...interface{}) (err error) {
	args = append([]interface{}{key}, args...)

	_, err = me.redisConn.Do("HMSET", args...)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[HMset]")
	}
	return
}

//HGetAll 返回hash值
func (me *Cache) HGetAll(key string, val interface{}) (err error) {
	v, err := redis.Values(me.redisConn.Do("HGETALL", key))
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
func (me *Cache) ZrangeInts(key string, start, stop int) (val []int) {
	var err error
	val, err = redis.Ints(me.redisConn.Do("ZRANGE", key, start, stop))
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[ZrangeInts]")
	}
	return
}

//ZrangeStrs 通过索引区间返回有序集合成指定区间内的成员
func (me *Cache) ZrangeStrs(key string, start, stop int) (val []string) {
	var err error
	val, err = redis.Strings(me.redisConn.Do("ZRANGE", key, start, stop))
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[ZrangeStrs]")
	}
	return
}

//LrangeInts 获取列表指定范围内的元素
func (me *Cache) LrangeInts(key string, start, stop int) (val []int) {
	var err error
	val, err = redis.Ints(me.redisConn.Do("LRANGE", key, start, stop))
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[LrangeInts]")
	}
	return
}

//LrangeStrs 获取列表指定范围内的元素
func (me *Cache) LrangeStrs(key string, start, stop int) (val []string) {
	var err error
	val, err = redis.Strings(me.redisConn.Do("LRANGE", key, start, stop))
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[LrangeStrs]")
	}
	return
}

//LremInt 移除列表中与参数 VALUE 相等的元素
func (me *Cache) LremInt(key string, count, value int) (err error) {
	_, err = me.redisConn.Do("LREM", key, count, value)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[LremInt]")
	}
	return
}

//RPush 向存于 key 的列表的尾部插入所有指定的值
func (me *Cache) RPush(key string, args []int) (err error) {
	values := redis.Args{}.Add(key).AddFlat(args)
	_, err = me.redisConn.Do("RPUSH", values...)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("[RPush]")
	}
	return
}
