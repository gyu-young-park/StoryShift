package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gyu-young-park/StoryShift/pkg/log"
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
func NewControllerManager(controllers ...controller) *controllerManager {
	logger := log.GetLogger()
	c := controllerManager{
		engine: gin.Default(),
	}

	apiGroupList := []string{}
	v1groupRouter := c.engine.Group("/v1")
	for _, apiController := range controllers {
		group := v1groupRouter.Group(apiController.GetAPIGroup())
		apiController.RegisterAPI(group)
		apiGroupList = append(apiGroupList, group.BasePath())
	}

	for _, apiGroup := range apiGroupList {
		logger.Info("API GROUP: " + apiGroup)
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
