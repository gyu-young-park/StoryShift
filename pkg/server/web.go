package server

import (
	"net/http"

	"github.com/gyu-young-park/VelogStoryShift/internal/config"
	"github.com/gyu-young-park/VelogStoryShift/pkg/controller"
	"github.com/gyu-young-park/VelogStoryShift/pkg/log"
)

func Start(c config.ConfigModel) {
	logger := log.GetLogger()
	s := http.Server{
		Addr:    ":8080",
		Handler: controller.Manager.GetHTTPHandler(),
	}

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Errorf("listen: %s\n", err)
	}
}
