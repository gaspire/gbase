package base

import (
	"os"
	"strconv"
)

//WechatType 微信类型
type WechatType int

// useropen 类型
const (
	WechatWeb  WechatType = iota // 0 网站
	WechatIOS                    // 1 ios
	WechatBiz                    // 2 微信公众号
	WechatMini                   // 3 小程序
)

//登录类型
const (
	LoginPassword  = 0 // 密码
	LoginCaptcha   = 1 // 验证码
	LoginMiniQuick = 2 // 微信快捷登录
)

// 验证码类型
const (
	SmsRegist        = 0 //注册
	SmsResetPassword = 1 //重置密码
	SmsChangeMobile  = 2 //变更手机号
	SmsLogin         = 3 //登录
	SmsCheck         = 4 //验证
)

const (
	//RoleStudent 学生
	RoleStudent = 2
	//RoleTeacher 外教
	RoleTeacher = 4
	//RoleAssistant 助教
	RoleAssistant = 8
)

// 支付方式
const (
	PayWayWechat     = 1 // 微信支付
	PayWayAlipay     = 2 // 支付宝
	PayWayBank       = 3 // 银联
	PayWayOffline    = 4 // 线下支付
	PayWayFree       = 5 // 免费
	PayWayWechatApp  = 6 // 小程序
	PayWayWechatOpen = 7 // 公众号
)

// 常量定义
const (
	DBConn    = "default,"
	RedisConn = "default"
)

var (
	OSS_HOST_CDN = os.Getenv("OSS_HOST_CDN")
	//IOSShowPay 是否显示支付,1-隐藏，2-显示
	IOSShowPay, _ = strconv.Atoi(os.Getenv("IOS_SHOW_PAY"))
	//RateLimit API频率限制
	RateLimit, _ = strconv.Atoi(os.Getenv("RateLimit"))
)
