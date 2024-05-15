package schema

type Register struct {
	Username   string `json:"username"`   // 用户名
	Password   string `json:"password"`   // 密码
	Repassword string `json:"repassword"` // 确认密码
}
