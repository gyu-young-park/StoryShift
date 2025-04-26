package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v1controller "github.com/gyu-young-park/VelogStoryShift/pkg/controller/v1"
)

var (
	v1Controllers = []controller{
		v1controller.NewhealthController("/"),
		v1controller.NewVelogController("/velog"),
	}
)

var Manager *controllerManager = newControllerManager()

type controller interface {
	RegisterAPI(*gin.RouterGroup)
	GetAPIGroup() string
}

type controllerManager struct {
	engine *gin.Engine
}

// TODO: change controller chaning like: c1.register(c2).register(c3)
// AND the path will be like c1/c2/c3/api
func newControllerManager() *controllerManager {
	c := controllerManager{
		engine: gin.Default(),
	}

	v1groupRouter := c.engine.Group("/v1")
	for _, apiController := range v1Controllers {
		group := v1groupRouter.Group(apiController.GetAPIGroup())
		apiController.RegisterAPI(group)
	}

	return &c
}

func (c controllerManager) GetHTTPHandler() http.Handler {
	return c.engine.Handler()
}
