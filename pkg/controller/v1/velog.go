package v1controller

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gyu-young-park/StoryShift/pkg/log"
	"github.com/gyu-young-park/StoryShift/pkg/service"
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
	router.GET("/post/download", downloadPost)
	router.GET("/posts", posts)
	router.POST("/posts/download", downloadSelectedPosts)
	router.GET("/posts/download", downloadAllPosts)
	router.GET("/series", series)
	router.GET("/series/download", downloadAllSeries)
	router.POST("/series/download", downloadSelectedSeries)
}

func post(c *gin.Context) {
	logger := log.GetLogger()
	var req VelogPostRequestModel
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Debugf("Velog username: %s, url_slug: %s",
		req.Name,
		req.URLSlug)

	velogPost, err := service.GetPost(req.Name, req.URLSlug)
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

	logger.Infof("Velog username: %s, post_id: %s, count: %v",
		req.Name,
		req.PostId,
		req.Count)

	velogPosts, err := service.GetPosts(req.Name, req.PostId, req.Count)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": velogPosts,
	})
}

func downloadPost(c *gin.Context) {
	logger := log.GetLogger()
	var req VelogPostRequestModel
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Debugf("Velog username: %s, url_slug: %s", req.Name, req.URLSlug)
	closeFunc, zipFilename, err := service.FetchVelogPostZip(req.Name, req.URLSlug)
	defer closeFunc()

	logger.Infof("get zip file: %s", zipFilename)
	if err != nil {
		logger.Errorf("failed to return zip file: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.FileAttachment(zipFilename, filepath.Base(zipFilename))
}

func downloadAllPosts(c *gin.Context) {
	logger := log.GetLogger()
	var req VelogUserNameReqestModel
	if err := c.ShouldBind(&req); err != nil {
		logger.Errorf("client error occureed: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	closeFunc, zipFilename, err := service.FetchAllVelogPostsZip(req.Name)
	defer closeFunc()
	if err != nil {
		logger.Errorf("server error occureed: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.FileAttachment(zipFilename, filepath.Base(zipFilename))
}

func downloadSelectedPosts(c *gin.Context) {
	var req []VelogPostRequestModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := ""
	urlSlugList := []string{}
	for _, data := range req {
		user = data.Name
		urlSlugList = append(urlSlugList, data.URLSlug)
	}

	closeFunc, zipFilename, err := service.FetchSelectedVelogPostsZip(user, urlSlugList)
	defer closeFunc()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.FileAttachment(zipFilename, filepath.Base(zipFilename))
}

func series(c *gin.Context) {
	c.String(http.StatusOK, "series")
}

func downloadAllSeries(c *gin.Context) {
	c.String(http.StatusOK, "all series")
}

func downloadSelectedSeries(c *gin.Context) {
	c.String(http.StatusOK, "all series")
}
