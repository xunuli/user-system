package constant

// 一些常量

// 一些常量，用作redis中key的值
const (
	ReqUuid          = "uuid"
	UserInfoPrefix   = "userinfo_"
	SessionKeyPrefix = "session_"
)

// 性别的限定值词
const (
	GenderMale   = "male"
	GenderFeMale = "female"
)

//设置session的key

const (
	SessionKey   = "user_session"
	CookieExpire = 3600
)
