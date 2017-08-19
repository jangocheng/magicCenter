package handler

import (
	"net/http"

	"muidea.com/magicCenter/application/common"
	"muidea.com/magicCenter/application/common/dbhelper"
	"muidea.com/magicCenter/application/common/model"
	"muidea.com/magicCenter/foundation/cache"
)

// CreateAuthorityHandler 新建CASHandler
func CreateAuthorityHandler(moduleHub common.ModuleHub, sessionRegistry common.SessionRegistry) common.AuthorityHandler {
	dbhelper, _ := dbhelper.NewHelper()
	i := impl{
		sessionRegistry:  sessionRegistry,
		authGroupManager: createAuthGroupManager(dbhelper),
		aclManager:       createACLManager(dbhelper),
		cacheData:        cache.NewCache()}

	casModule, _ := moduleHub.FindModule(common.CASModuleID)
	entryPoint := casModule.EntryPoint()
	switch entryPoint.(type) {
	case common.CASHandler:
		i.casHandler = entryPoint.(common.CASHandler)
	}

	return &i
}

type impl struct {
	sessionRegistry  common.SessionRegistry
	casHandler       common.CASHandler
	authGroupManager authGroupManager
	aclManager       aclManager
	cacheData        cache.Cache
}

/*
1、先判断authToken是否一致，如果不一致则，认为无权限
*/
func (i *impl) VerifyAuthority(res http.ResponseWriter, req *http.Request) bool {
	session := i.sessionRegistry.GetSession(res, req)

	_, ok := session.GetAccount()
	if !ok {
		return false
	}

	urlToken := req.URL.Query().Get(common.AuthTokenID)
	sessionToken := ""
	obj, ok := session.GetOption(common.AuthTokenID)
	if ok {
		// 找到sessionToken了，则说明该用户已经登录了，这里就必须保证两端的token一致否则也要认为鉴权非法
		// 用户登录过Token必然不为空
		sessionToken = obj.(string)

		if i.casHandler != nil {
			i.casHandler.RefreshToken(sessionToken, req.RemoteAddr)
		}

		return urlToken == sessionToken
	}

	return true
}

func (i *impl) QueryAuthGroup(module string) ([]model.AuthGroup, bool) {
	return i.authGroupManager.queryAuthGroup(module)
}

func (i *impl) InsertAuthGroup(authGroups []model.AuthGroup) bool {
	return i.authGroupManager.insertAuthGroup(authGroups)
}

func (i *impl) DeleteAuthGroup(authGroups []model.AuthGroup) bool {
	return i.authGroupManager.deleteAuthGroup(authGroups)
}

func (i *impl) AdjustUserAuthGroup(userID int, authGroup []int) bool {
	return i.authGroupManager.adjustUserAuthGroup(userID, authGroup)
}

func (i *impl) GetUserAuthGroup(userID int) ([]int, bool) {
	return i.authGroupManager.getUserAuthGroup(userID)
}

func (i *impl) QueryACL(module string, status int) ([]model.ACL, bool) {
	return i.aclManager.queryACL(module, status)
}

func (i *impl) AddACL(url, method, module string) (model.ACL, bool) {
	return i.aclManager.addACL(url, method, module)
}

func (i *impl) DelACL(url, method, module string) bool {
	return i.aclManager.delACL(url, method, module)
}

func (i *impl) EnableACL(acls []int) bool {
	return i.aclManager.enableACL(acls)
}

func (i *impl) DisableACL(acls []int) bool {
	return i.aclManager.disableACL(acls)
}

func (i *impl) AdjustACLAuthGroup(acl int, authGroup []int) (model.ACL, bool) {
	return i.aclManager.adjustACLAuthGroup(acl, authGroup)
}
