package schema

type Register struct {
	Nickname   string `json:"nickname"`   // 昵称
	Username   string `json:"username"`   // 用户名
	Password   string `json:"password"`   // 密码
	Repassword string `json:"repassword"` // 确认密码
	Avatar     string `json:"avatar"`     // 头像
}
