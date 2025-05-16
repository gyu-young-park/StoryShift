package v1statuscontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	servicestatus "github.com/gyu-young-park/StoryShift/pkg/service/status"
)

type statueController struct {
	APIGroup string
}

func NewStatueController(apiGroup string) *statueController {
	return &statueController{
		APIGroup: apiGroup,
	}
}

func (h *statueController) GetAPIGroup() string {
	return h.APIGroup
}

func (v *statueController) RegisterAPI(router *gin.RouterGroup) {
	router.GET("/startup", startupCheckHandler)
	router.GET("/health", livenessCheckHandler)
	router.GET("/ready", readinessCheckHandler)
}

func startupCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func livenessCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "health"})
}

func readinessCheckHandler(c *gin.Context) {
	if !servicestatus.Ready() {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "server is not ready"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ready": "ready"})
}
