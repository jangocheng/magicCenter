package application

import (
	"os"
	"webcenter/kernel"
	_"webcenter/kernel/admin"
	_"webcenter/module/loader"
)

var serverPort string = "8888"

type Application interface {
	Run()
}

var app *application = nil

type application struct {
}

func AppInstance() Application {
	if (app == nil) {
		app = &application{}
		
		app.construct()
	}
	
	return app 
}

func (instance application) construct() {
	os.Setenv("PORT", serverPort)
}

func (instance application) Run() {
	kernel.Initialize()
	
	kernel.Run()
	
	kernel.Uninitialize()	
}


