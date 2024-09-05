package schema

type RegisterData struct {
	Username   string `json:"username"`   // 用户名
	Password   string `json:"password"`   // 密码
	Repassword string `json:"repassword"` // 确认密码
	Nickname   string `json:"nickname"`   // 昵称
	Avatar     string `json:"avatar"`     // 头像
}

type BindData struct {
	Uid        uint64 `json:"uid"`        // 设备名称
	Username   string `json:"username"`   // 用户名
	Password   string `json:"password"`   // 密码
	Repassword string `json:"repassword"` // 确认密码
}
