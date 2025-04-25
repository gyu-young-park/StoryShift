package web

import "github.com/gyu-young-park/VelogStoryShift/internal/config"

type server interface {
	Run(config.ConfigModel) error
	Close() error
}

func Server(c config.ConfigModel) server {
	return newGinServer(c)
}
