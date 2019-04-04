package base

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Controller 基类
type Controller struct {
}

//Success 统一输出正确结果
func (me *Controller) Success(c *gin.Context, output interface{}) {
	data := SuccessResponse{
		Code:       200,
		Exp:        c.GetInt64("exp"),
		ServerTime: time.Now().Format("2006-01-02 15:04:05"),
		Tip:        "success",
		Output:     output,
	}
	c.JSON(200, data)
}

//Error 统一输出错误结果
func (me *Controller) Error(c *gin.Context, code int) {
	key := strconv.Itoa(code)
	data := ErrResponse{
		Code:       code,
		Exp:        c.GetInt64("exp"),
		ServerTime: time.Now().Format("2006-01-02 15:04:05"),
		Tip:        ErrCode[key],
	}
	c.JSON(200, data)
}
