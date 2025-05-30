package v1statuscontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	servicestatus "github.com/gyu-young-park/StoryShift/pkg/service/status"
)

type StatueController struct {
	service  *servicestatus.StatusService
	APIGroup string
}

func NewStatueController(service *servicestatus.StatusService) *StatueController {
	return &StatueController{
		service:  service,
		APIGroup: "/status",
	}
}

func (s *StatueController) GetAPIGroup() string {
	return s.APIGroup
}

func (s *StatueController) RegisterAPI(router *gin.RouterGroup) {
	router.GET("/startup", s.startupCheckHandler)
	router.GET("/health", s.livenessCheckHandler)
	router.GET("/ready", s.readinessCheckHandler)
}

func (s *StatueController) startupCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *StatueController) livenessCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "health"})
}

func (s *StatueController) readinessCheckHandler(c *gin.Context) {
	if !s.service.Ready() {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "server is not ready"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ready": "ready"})
}
