package main

import (
	"github.com/gyu-young-park/StoryShift/internal/config"
	"github.com/gyu-young-park/StoryShift/pkg/log"
	"github.com/gyu-young-park/StoryShift/pkg/server"
	"go.uber.org/zap"
)

func main() {
	logger := log.GetLogger()
	logger.Info("App starts", zap.String("hellp", "world"))
	server.Start(config.Manager.ConfigModel)
}
