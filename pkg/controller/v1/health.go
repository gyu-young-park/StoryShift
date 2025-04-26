package v1controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type healthController struct {
	APIGroup string
}

func NewhealthController(apiGroup string) *healthController {
	return &healthController{
		APIGroup: apiGroup,
	}
}

func (h *healthController) GetAPIGroup() string {
	return h.APIGroup
}

func (v *healthController) RegisterAPI(router *gin.RouterGroup) {
	router.GET("health", health)
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}
