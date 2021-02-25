package base

import (
	"time"

	"github.com/gin-gonic/gin"
)

//ErrCode 错误码
var ErrCode map[int]string

//Response 返回数据结构
type Response struct {
	Code       int         `json:"code"`
	IOSShowPay int         `json:"ios_pay,omitempty"`
	Exp        int64       `json:"exp,omitempty"`
	ServerTime string      `json:"server_time"`
	Tip        string      `json:"tip"`
	Output     interface{} `json:"output,omitempty"`
}

// Controller 基类
type Controller struct {
}

//Success 统一输出正确结果
func (me *Controller) Success(c *gin.Context, output interface{}) {
	data := Response{
		Code:       200,
		Exp:        c.GetInt64("exp"),
		IOSShowPay: IOSShowPay,
		ServerTime: time.Now().Format("2006-01-02 15:04:05"),
		Tip:        "success",
		Output:     output,
	}
	c.JSON(200, data)
}

//Error 统一输出错误结果
func (me *Controller) Error(c *gin.Context, code int) {
	data := Response{
		Code:       code,
		Exp:        c.GetInt64("exp"),
		ServerTime: time.Now().Format("2006-01-02 15:04:05"),
		Tip:        ErrCode[code],
	}
	c.JSON(200, data)
}
