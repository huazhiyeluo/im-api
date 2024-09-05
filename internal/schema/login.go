package schema

type LoginData struct {
	Platform string `json:"platform"` // 平台      visitor | account(用户名|手机号|邮箱) | facebook
	Username string `json:"username"` // 用户名
	Phone    string `json:"phone"`    // 手机号
	Email    string `json:"email"`    // 邮箱
	Password string `json:"password"` // 密码
	Nickname string `json:"nickname"` // 昵称
	Avatar   string `json:"avatar"`   // 头像
	Token    string `json:"token"`    // 1、token登录的token | 2、fb的access_token
	Siteuid  string `json:"siteuid"`  // 1、fb的openid
}

type PublicVar struct {
	IsNewUser uint32
	Siteuid   string
	Sid       uint32
	Uid       uint64
	Nickname  string
	Avatar    string
	Phone     string
	Email     string
	Err       error
}
