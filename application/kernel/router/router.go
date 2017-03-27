package router

import (
	"log"
	"net/http"

	"muidea.com/magicCenter/application/common"

	"github.com/go-martini/martini"
)

// Router 路由器对象
type Router interface {
	// 增加路由
	AddRoute(rt common.Route)
	// 清除路由
	RemoveRoute(rt common.Route)

	// 增加Get路由
	AddGetRoute(pattern string, handler interface{})
	// 清除Get路由
	RemoveGetRoute(pattern string)
	// 增加Post路由
	AddPostRoute(pattern string, handler interface{})
	// 清除Post路由
	RemovePostRoute(pattern string)
	// 增加Delete路由
	AddDeleteRoute(pattern string, handler interface{})
	// 清除Delete路由
	RemoveDeleteRoute(pattern string)
	// 增加Put路由
	AddPutRoute(pattern string, handler interface{})
	// 清除Put路由
	RemovePutRoute(pattern string)

	// 返回Martini.Router对象
	Router() martini.Router

	// 分发一条请求
	Dispatch(res http.ResponseWriter, req *http.Request)
}

// CreateRouter 新建Router
func CreateRouter() Router {
	impl := impl{}
	impl.martiniRouter = martini.NewRouter()

	return &impl
}

// impl 路由器实现
type impl struct {
	martiniRouter martini.Router
}

// AddRoute 增加Route
func (instance *impl) AddRoute(rt common.Route) {
	switch rt.Method() {
	case common.GET:
		instance.AddGetRoute(rt.Pattern(), rt.Handler())
	case common.POST:
		instance.AddPostRoute(rt.Pattern(), rt.Handler())
	case common.DELETE:
		instance.AddDeleteRoute(rt.Pattern(), rt.Handler())
	case common.PUT:
		instance.AddPutRoute(rt.Pattern(), rt.Handler())
	}
}

// RemoveRoute 清除Route
func (instance *impl) RemoveRoute(rt common.Route) {
	switch rt.Method() {
	case common.GET:
		instance.RemoveGetRoute(rt.Pattern())
	case common.POST:
		instance.RemovePostRoute(rt.Pattern())
	case common.DELETE:
		instance.RemoveDeleteRoute(rt.Pattern())
	case common.PUT:
		instance.RemovePutRoute(rt.Pattern())
	}
}

// Router 返回系统Router
func (instance *impl) Router() martini.Router {
	return instance.martiniRouter
}

// AddGetRoute 添加一条Get路由
func (instance *impl) AddGetRoute(pattern string, handler interface{}) {
	if martini.Env != martini.Prod {
		log.Printf("[get]:%s", pattern)
	}

	instance.martiniRouter.Get(pattern, handler)
}

// RemoveGetRoute 清除一条Get路由
func (instance *impl) RemoveGetRoute(pattern string) {
}

// AddPutRoute 添加一条Put路由
func (instance *impl) AddPutRoute(pattern string, handler interface{}) {
	if martini.Env != martini.Prod {
		log.Printf("[put]:%s", pattern)
	}

	instance.martiniRouter.Put(pattern, handler)
}

// RemovePutRoute 清除一条Put路由
func (instance *impl) RemovePutRoute(pattern string) {
}

// AddPostRoute 添加一条Post路由
func (instance *impl) AddPostRoute(pattern string, handler interface{}) {
	if martini.Env != martini.Prod {
		log.Printf("[post]:%s", pattern)
	}

	instance.martiniRouter.Post(pattern, handler)
}

// RemovePostRoute 清除一条Post路由
func (instance *impl) RemovePostRoute(pattern string) {
}

// AddDeleteRoute 添加一条Delete路由
func (instance *impl) AddDeleteRoute(pattern string, handler interface{}) {
	if martini.Env != martini.Prod {
		log.Printf("[delete]:%s", pattern)
	}

	instance.martiniRouter.Delete(pattern, handler)
}

// RemoveDeleteRoute 清除一条Delete路由
func (instance *impl) RemoveDeleteRoute(pattern string) {
}

// Dispatch 分发一次请求
func (instance *impl) Dispatch(res http.ResponseWriter, req *http.Request) {

}
