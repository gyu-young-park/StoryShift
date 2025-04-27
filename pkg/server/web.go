package server

import (
	"fmt"
	"net/http"

	"github.com/gyu-young-park/VelogStoryShift/internal/config"
	"github.com/gyu-young-park/VelogStoryShift/pkg/controller"
	"github.com/gyu-young-park/VelogStoryShift/pkg/log"
)

func Start(c config.ConfigModel) {
	logger := log.GetLogger()

	s := http.Server{
		Addr:    fmt.Sprintf(":%s", config.Manager.AppConfig.Server.Port),
		Handler: controller.Manager.GetHTTPHandler(),
	}

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Errorf("listen: %s\n", err)
	}
}
