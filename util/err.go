package util

import (
	log "github.com/sirupsen/logrus"
)

// CommonError 微信返回的通用错误json
type CommonError struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// FailOnError log error
func FailOnError(err error, msg string) bool {
	if err != nil {
		log.WithFields(log.Fields{"err": err, "msg": msg}).Error("FailOnError")
		return true
	}
	return false
}
