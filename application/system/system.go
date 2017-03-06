package system

import (
	"net/http"

	"muidea.com/magicCenter/application/common"
	"muidea.com/magicCenter/application/common/service"
	"muidea.com/magicCenter/application/system/authority"
	"muidea.com/magicCenter/application/system/modulehub"
	"muidea.com/magicCenter/application/system/router"
	"muidea.com/magicCenter/application/system/session"

	"github.com/go-martini/martini"
)

type impl struct {
	loaderImpl          service.ModuleLoader
	configurationImpl   service.Configuration
	routerImpl          service.Router
	moduleHubImpl       service.ModuleHub
	authorityImpl       service.Authority
	sessionRegistryImpl service.SessionRegistry
	instanceFrameImpl   *martini.Martini
}

// NewSystem 新建System对象
func NewSystem(loader service.ModuleLoader, configuration service.Configuration) service.System {
	i := &impl{
		loaderImpl:          loader,
		configurationImpl:   configuration,
		routerImpl:          router.CreateRouter(),
		moduleHubImpl:       modulehub.CreateModuleHub(),
		authorityImpl:       authority.CreateAuthority(),
		sessionRegistryImpl: session.CreateSessionRegistry(),
		instanceFrameImpl:   martini.New()}

	return i
}

func (i *impl) StartUp() error {
	i.configurationImpl.LoadConfig()

	i.loaderImpl.LoadAllModules(i)

	allModules := i.moduleHubImpl.QueryAllModule()
	for _, m := range allModules {
		baseURL := m.URL()
		routes := m.Routes()
		for _, rt := range routes {
			i.routerImpl.AddRoute(baseURL, rt)
		}
	}

	i.moduleHubImpl.StartupAllModules()
	return nil
}

func (i *impl) Run() {
	martiniRouter := i.routerImpl.Router()

	i.instanceFrameImpl.Use(martini.Logger())
	i.instanceFrameImpl.Use(martini.Recovery())
	i.instanceFrameImpl.Use(authority.Authority(i))

	i.instanceFrameImpl.MapTo(martiniRouter, (*martini.Routes)(nil))
	i.instanceFrameImpl.Action(martiniRouter.Handle)

	instance := martini.ClassicMartini{}
	instance.Martini = i.instanceFrameImpl
	instance.Router = martiniRouter

	instance.Run()
}

func (i *impl) ShutDown() error {
	/*
		退出时不需要做路由清理操作
		allModules := i.moduleHubImpl.QueryAllModule()
		for _, m := range allModules {
			baseURL := m.URL()
			routes := m.Routes()
			for _, rt := range routes {
				routerImpl.RemoveRoute(baseURL, rt)
			}
		}
	*/

	i.moduleHubImpl.CleanupAllModules()
	return nil
}

// GetRouter 获取系统的Router
func (i *impl) Router() service.Router {
	return i.routerImpl
}

// GetModuleHub 获取系统的ModuleHub
func (i *impl) ModuleHub() service.ModuleHub {
	return i.moduleHubImpl
}

// GetConfiguration 获取当前Configuration
func (i *impl) Configuration() service.Configuration {
	return i.configurationImpl
}

func (i *impl) Authority() service.Authority {
	return i.authorityImpl
}

// GetSession 获取当前Session
func (i *impl) Session(w http.ResponseWriter, r *http.Request) common.Session {
	return i.sessionRegistryImpl.GetSession(w, r)
}
