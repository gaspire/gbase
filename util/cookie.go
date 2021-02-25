package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//CookieHandle 用来做cookie 处理
type CookieHandle struct {
	HTTPWriter  http.ResponseWriter //主要用来写入头部
	HTTPRequest *http.Request       //主要用来获取头部信息
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	HTTPOnly bool
	Domain   string
	Path     string
}

//CookieInit 初始化
func CookieInit(w gin.ResponseWriter, req *http.Request, expire int, HTTPOnly bool) *CookieHandle {
	//cookie 初始化
	cookie := new(CookieHandle)
	cookie.HTTPWriter = w
	cookie.HTTPRequest = req
	cookie.MaxAge = expire
	cookie.HTTPOnly = HTTPOnly
	cookie.Domain = "waypal.com"
	cookie.Path = "/"
	return cookie
}

//GetCookie 获取cookie
func (me *CookieHandle) GetCookie(key string) string {
	cookie, err := me.HTTPRequest.Cookie(key)
	if err != nil {
		return ""
	}
	return cookie.Value
}

//SetCookie 设置cookie param[0]=>过期时间
func (me *CookieHandle) SetCookie(key, val string, expire ...int) {
	//创建 http.cookie 指针
	cookie := new(http.Cookie)
	cookie.Name = key
	cookie.Value = val
	cookie.Domain = "waypal.com"
	cookie.Path = "/"
	cookie.HttpOnly = me.HTTPOnly
	//如果有传递超时时间 重新设置超时时间
	if len(expire) > 0 {
		me.MaxAge = expire[0]
	}
	//如果不设置 过期时间 不赋值
	if me.MaxAge != 0 {
		cookie.MaxAge = me.MaxAge
	}
	//调用http 包 进行设置cookie  大致是 讲cookie 构造体的数据 生成 string  header().set()
	http.SetCookie(me.HTTPWriter, cookie)
}

//DelCookie 删除cookie
func (me *CookieHandle) DelCookie(key string) {
	cookie := new(http.Cookie)
	cookie.Name = key
	//讲过期时间改为-1就可以清空对应的cookie
	cookie.MaxAge = -1
	http.SetCookie(me.HTTPWriter, cookie)
}
