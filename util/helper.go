package util

//获取随机数
import (
	_ "image/jpeg"
	_ "image/png"

	"github.com/mojocn/base64Captcha"
)

//GenerateCharacterCaptcha 生成图片字符验证码
func GenerateCharacterCaptcha(height, width, mode, len int) (captchaID, base64Png string) {
	config := base64Captcha.ConfigCharacter{
		Height:             height,
		Width:              width,
		Mode:               mode,
		ComplexOfNoiseText: 0,
		ComplexOfNoiseDot:  0,
		IsUseSimpleFont:    true,
		IsShowHollowLine:   false,
		IsShowNoiseDot:     false,
		IsShowNoiseText:    false,
		IsShowSlimeLine:    false,
		IsShowSineLine:     false,
		CaptchaLen:         len,
	}
	captchaID, digitCap := base64Captcha.GenerateCaptcha("", config)
	base64Png = base64Captcha.CaptchaWriteToBase64Encoding(digitCap)
	return
}

//VerfiyCaptcha 验证
func VerfiyCaptcha(captchaID, captcha string) bool {
	return base64Captcha.VerifyCaptcha(captchaID, captcha)
}
