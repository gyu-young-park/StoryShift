package web

import (
	"github.com/gin-gonic/gin"
	"github.com/gyu-young-park/VelogStoryShift/internal/config"
)

func newGinServer(c config.ConfigModel) *ginServer {
	return &ginServer{
		engine: gin.Default(),
	}
}

type ginServer struct {
	engine *gin.Engine
}

func (g *ginServer) Run(c config.ConfigModel) error {
	g.engine.GET("hello", func(con *gin.Context) {
		con.JSON(200, gin.H{
			"message": "pong",
		})
	})
	return g.engine.Run()
}

func (g *ginServer) Close() error {
	return nil
}
