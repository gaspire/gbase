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
	// SubmissionFee 摘要费用
	SubmissionFee, _ = strconv.ParseFloat(os.Getenv("SUBMISSION_FEE"), 64)
	// ApplyRegularFee 普通用户费用
	ApplyRegularFee, _ = strconv.ParseFloat(os.Getenv("APPLY_REGULAR_FEE"), 64)
	// ApplyStudentFee 学生费用
	ApplyStudentFee, _ = strconv.ParseFloat(os.Getenv("APPLY_STUDENT_FEE"), 64)
)
