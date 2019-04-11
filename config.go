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

const (
	ROLE_REGISTER = 2 //普通注册用户
)

// 支付方式
const (
	PayPalWay = 1 // paypal支付
)

var (
	//RateLimit API频率限制
	RateLimit, _ = strconv.Atoi(os.Getenv("RateLimit"))
)
