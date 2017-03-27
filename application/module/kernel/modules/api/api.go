package api

import (
	"muidea.com/magicCenter/application/common"
	"muidea.com/magicCenter/application/common/model"
	"muidea.com/magicCenter/application/module/kernel/modules/api/def"
	"muidea.com/magicCenter/application/module/kernel/modules/api/route"
)

type api struct {
	moduleHub       common.ModuleHub
	sessionRegistry common.SessionRegistry
	routes          []common.Route
}

// LoadModule 加载模块
func LoadModule(cfg common.Configuration, sessionRegistry common.SessionRegistry, modHub common.ModuleHub) {
	instance := &api{moduleHub: modHub, sessionRegistry: sessionRegistry, routes: []common.Route{}}
	instance.routes = route.AppendArticleRoute(instance.routes, modHub, sessionRegistry)
	instance.routes = route.AppendCatalogRoute(instance.routes, modHub, sessionRegistry)
	instance.routes = route.AppendLinkRoute(instance.routes, modHub, sessionRegistry)
	instance.routes = route.AppendMediaRoute(instance.routes, modHub, sessionRegistry)
	instance.routes = route.AppendUserRoute(instance.routes, modHub)
	instance.routes = route.AppendGroupRoute(instance.routes, modHub)

	modHub.RegisterModule(instance)
}

func (instance *api) ID() string {
	return def.ID
}

func (instance *api) Name() string {
	return def.Name
}

func (instance *api) Description() string {
	return def.Description
}

func (instance *api) Group() string {
	return "admin api"
}

func (instance *api) Type() int {
	return common.KERNEL
}

func (instance *api) Status() int {
	return 0
}

func (instance *api) EndPoint() interface{} {
	return nil
}

func (instance *api) AuthGroups() []model.AuthGroup {
	groups := []model.AuthGroup{}

	groups = append(groups, model.CreateAuthGroup("PublicGroup", "允许查看公开权限的内容", def.ID))
	groups = append(groups, model.CreateAuthGroup("UserGroup", "允许查看用户权限范围内的内容", def.ID))
	groups = append(groups, model.CreateAuthGroup("AdminGroup", "允许管理用户权限范围内的内容", def.ID))

	return groups
}

// Route 路由信息
func (instance *api) Routes() []common.Route {
	return instance.routes
}

// Startup 启动模块
func (instance *api) Startup() bool {
	return true
}

// Cleanup 清除模块
func (instance *api) Cleanup() {

}
