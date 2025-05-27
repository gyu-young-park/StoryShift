package server

import (
	"fmt"
	"net/http"

	"github.com/gyu-young-park/StoryShift/internal/config"
	"github.com/gyu-young-park/StoryShift/pkg/controller"
	"github.com/gyu-young-park/StoryShift/pkg/log"
)

func Start(c config.ConfigModel) {
	logger := log.GetLogger()

	s := http.Server{
		Addr:    fmt.Sprintf(":%s", config.Manager.AppConfig.Server.Port),
		Handler: controller.NewControllerManager().GetHTTPHandler(),
	}
	logger.Infof("Server started, port: %v", config.Manager.AppConfig.Server.Port)

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Errorf("listen: %s\n", err)
	}
}
