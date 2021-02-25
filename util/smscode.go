package util

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	log "github.com/sirupsen/logrus"
)

//TemplateCode 短信模版类型
var TemplateCode = map[int]string{
	0:  "SMS_122810010", //注册
	1:  "SMS_122810009", //重置密码
	2:  "SMS_122810008", //变更手机号
	3:  "SMS_122810012", //登录
	4:  "SMS_122810014", //身份验证
	10: "SMS_136176349", //上课提醒
}

// SendSms 发送短信
func SendSms(mobile string, params []byte, msgType int) error {
	templateCode, ok := TemplateCode[msgType]
	if !ok {
		templateCode = TemplateCode[0]
	}

	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", "LTAIVEC3rPRRES79", "vm2gBV70XfY9To5GDGJ4zzfahrDafh")
	if err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("SendSms")
		return err
	}

	request := requests.NewCommonRequest()
	request.Method = "GET"
	request.Scheme = "http" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-hangzhou"
	request.QueryParams["PhoneNumbers"] = mobile
	request.QueryParams["SignName"] = "WayPal英语"
	request.QueryParams["TemplateCode"] = templateCode
	if params != nil {
		request.QueryParams["TemplateParam"] = string(params)
	}
	_, err = client.ProcessCommonRequest(request)
	if err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("SendSms")
	}
	return err
}
