package common

import (
	"muidea.com/magicCenter/application/common/model"
)

// CASHandler 鉴权处理Handler
type CASHandler interface {
	//@in account 账号
	//@in password 密码
	//@in remoteAddr 登陆地址
	//@ret model.UserDetail 登陆用户
	//@ret string 本次登陆的鉴权token
	//@ret bool 是否登陆成功
	LoginAccount(account, password, remoteAddr string) (model.UserDetail, string, bool)

	LoginToken(authToken, remoteAddr string) (string, bool)

	//@in authToken 鉴权token
	//@ret bool 是否登出成功
	Logout(authToken, remoteAddr string) bool

	// 刷新Token
	RefreshToken(authToken string, remoteAddr string) bool

	// 校验权限是否OK
	VerifyToken(authToken string) (model.OnlineAccountInfo, bool)

	// 生成一个静态Token
	AllocStaticToken(id string, expiration int64) (string, bool)
}
