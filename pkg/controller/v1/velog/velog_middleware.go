package v1velogcontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	servicevelog "github.com/gyu-young-park/StoryShift/pkg/service/velog"
)

func validateVelogUser(service *servicevelog.VelogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.Param("user")
		if !service.IsVelogUserExists(user) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Can't find the user: " + user,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
