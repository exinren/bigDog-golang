package errcode

var (
	ErrorGetCaptchaFail = NewError(20010001, "获取验证码失败")
	ErrorVerityCaptchaFail = NewError(20010002, "验证码验证失败")
)
