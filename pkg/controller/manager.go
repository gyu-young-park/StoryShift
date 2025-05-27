package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gyu-young-park/StoryShift/internal/injector"
	v1statuscontroller "github.com/gyu-young-park/StoryShift/pkg/controller/v1/status"
	v1velogcontroller "github.com/gyu-young-park/StoryShift/pkg/controller/v1/velog"
)

type controller interface {
	RegisterAPI(*gin.RouterGroup)
	GetAPIGroup() string
}

type controllerManager struct {
	engine *gin.Engine
}

// TODO: change controller chaning like: c1.register(c2).register(c3)
// AND the path will be like c1/c2/c3/api
func NewControllerManager() *controllerManager {
	c := controllerManager{
		engine: gin.Default(),
	}

	v1Controllers := []controller{
		v1statuscontroller.NewStatueController("/status", injector.Container.StatusService()),
		v1velogcontroller.NewVelogController("/velog", injector.Container.VelogService()),
	}

	apiGroupList := []string{}
	v1groupRouter := c.engine.Group("/v1")
	for _, apiController := range v1Controllers {
		group := v1groupRouter.Group(apiController.GetAPIGroup())
		apiController.RegisterAPI(group)
		apiGroupList = append(apiGroupList, group.BasePath())
	}

	c.engine.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"api": apiGroupList,
		})
	})

	return &c
}

func (c controllerManager) GetHTTPHandler() http.Handler {
	return c.engine.Handler()
}
