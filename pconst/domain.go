package pconst

import "errors"

// Domains in geektime site
const (
	GeekBang             = "https://time.geekbang.org"
	GeekBangAccount      = "https://account.geekbang.org"
	GeekBangCookieDomain = ".geekbang.org"
)

const (
	// UserAgentHeaderName ...
	UserAgentHeaderName = "User-Agent"
	// OriginHeaderName ...
	OriginHeaderName = "Origin"
	// UserAgentHeaderValue ...
	UserAgentHeaderValue = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.92 Safari/537.36"
	//UserAgentHeaderValue = "Mozilla/5.0 (Linux; U; Android 11; zh-CN; HD1900 Build/RKQ1.201022.002) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/78.0.3904.108 Quark/5.8.0.220 Mobile Safari/537.36"
	//UserAgentHeaderValue = "Dalvik/2.1.0 (Linux; U; Android 11; HD1900 Build/RKQ1.201022.002)"
)

var (
	// ErrGeekTimeRateLimit ...
	ErrGeekTimeRateLimit = errors.New("已触发限流, 你可以选择重新登录/重新获取 cookie, 或者稍后再试, 然后生成剩余的文章")
	// ErrAuthFailed ...
	ErrAuthFailed = errors.New("当前账户在其他设备登录或者登录已经过期, 请尝试重新登录")
)

var (
	// ErrWrongPassword ...
	ErrWrongPassword = errors.New("密码错误, 请尝试重新登录")
	// ErrTooManyLoginAttemptTimes ...
	ErrTooManyLoginAttemptTimes = errors.New("密码输入错误次数过多，已触发验证码校验，请稍后再试")
)
