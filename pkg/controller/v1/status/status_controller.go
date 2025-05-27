package v1statuscontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	servicestatus "github.com/gyu-young-park/StoryShift/pkg/service/status"
)

type statueController struct {
	service  *servicestatus.StatusService
	APIGroup string
}

func NewStatueController(apiGroup string, service *servicestatus.StatusService) *statueController {
	return &statueController{
		service:  service,
		APIGroup: apiGroup,
	}
}

func (s *statueController) GetAPIGroup() string {
	return s.APIGroup
}

func (s *statueController) RegisterAPI(router *gin.RouterGroup) {
	router.GET("/startup", s.startupCheckHandler)
	router.GET("/health", s.livenessCheckHandler)
	router.GET("/ready", s.readinessCheckHandler)
}

func (s *statueController) startupCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *statueController) livenessCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "health"})
}

func (s *statueController) readinessCheckHandler(c *gin.Context) {
	if !s.service.Ready() {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "server is not ready"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ready": "ready"})
}
