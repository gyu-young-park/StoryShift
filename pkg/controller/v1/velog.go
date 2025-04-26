package v1controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gyu-young-park/VelogStoryShift/internal/config"
	"github.com/gyu-young-park/VelogStoryShift/pkg/log"
	"github.com/gyu-young-park/VelogStoryShift/pkg/velog"
)

type velogController struct {
	APIGroup string
}

func NewVelogController(apiGroup string) *velogController {
	return &velogController{
		APIGroup: apiGroup,
	}
}

func (v *velogController) GetAPIGroup() string {
	return v.APIGroup
}

func (v *velogController) RegisterAPI(router *gin.RouterGroup) {
	router.GET("/post", post)
	router.GET("/posts", posts)
}

func post(c *gin.Context) {
	logger := log.GetLogger()
	var req VelogPostRequestModel
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Debugf("Velog username: %s, url_slog: %s",
		req.Name,
		req.URLSlog)
	velogApi := velog.NewVelogAPI(config.Manager.VelogConfig.URL, req.Name)
	velogPost, err := velogApi.Post(req.URLSlog)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": velogPost,
	})
}

func posts(c *gin.Context) {
	logger := log.GetLogger()
	var req VelogPostsRequestModel
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Infof("Velog username: %s, url_slog: %s, limit: %v",
		req.Name,
		req.PostId,
		req.Limit)
	velogApi := velog.NewVelogAPI(config.Manager.VelogConfig.URL, req.Name)
	velogPosts, err := velogApi.Posts(req.PostId, req.Limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": velogPosts,
	})
}
