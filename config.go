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
	PayPalWay = 4 // paypal支付
)

var (
	//RateLimit API频率限制
	RateLimit, _ = strconv.Atoi(os.Getenv("RateLimit"))
	// ApplyFee 报名费用/单价
	ApplyFee, _ = strconv.ParseFloat(os.Getenv("APPLY_FEE"), 64)
)
