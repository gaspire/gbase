package base

import (
	"os"
	"strconv"
)

// 常量定义
const (
	DBConn    = "default,"
	RedisConn = "default"
)

var (
	//RateLimit API频率限制
	RateLimit, _ = strconv.Atoi(os.Getenv("RateLimit"))
)
