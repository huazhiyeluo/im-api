package login

import (
	"qqapi/internal/schema"
)

// FacebookLogin Facebook 登录适配器
type FacebookLogin struct {
	cin *schema.CommonData
	in  *schema.LoginData
	pv  *schema.PublicVar
}

func (b *FacebookLogin) Verify() {

}

func (b *FacebookLogin) IsNewUser() {
}
