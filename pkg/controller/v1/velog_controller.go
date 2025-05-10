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
	router.GET("/:user/post", getPost)
	router.GET("/:user/post/download", downloadPost)
	router.GET("/:user/posts", getPosts)
	router.POST("/:user/posts/download", downloadSelectedPosts)
	router.GET("/:user/posts/download", downloadAllPosts)
	router.GET("/:user/series", getSeries)
	router.GET("/:user/series/:url_slug", getPostsInSeries)
	router.GET("/:user/series/download", downloadAllSeries)
	router.POST("/:user/series/download", downloadSelectedSeries)
}

func getPost(c *gin.Context) {
	logger := log.GetLogger()
	user := c.Param("user")
	var req VelogPostRequestModel
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Debugf("Velog username: %s, url_slug: %s",
		user,
		req.URLSlug)

	velogPost, err := service.GetPost(user, req.URLSlug)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": velogPost,
	})
}

func getPosts(c *gin.Context) {
	logger := log.GetLogger()
	user := c.Param("user")

	var req VelogPostsRequestModel
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Infof("Velog username: %s, post_id: %s, count: %v",
		user,
		req.PostId,
		req.Count)

	velogPosts, err := service.GetPosts(user, req.PostId, req.Count)
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
	user := c.Param("user")

	var req VelogPostRequestModel
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Debugf("Velog username: %s, url_slug: %s", user, req.URLSlug)
	closeFunc, zipFilename, err := service.FetchVelogPostZip(user, req.URLSlug)
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
	user := c.Param("user")
	// var req VelogUserNameReqestModel
	// if err := c.ShouldBind(&req); err != nil {
	// 	logger.Errorf("client error occureed: %s", err)
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	closeFunc, zipFilename, err := service.FetchAllVelogPostsZip(user)
	defer closeFunc()
	if err != nil {
		logger.Errorf("server error occureed: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.FileAttachment(zipFilename, filepath.Base(zipFilename))
}

func downloadSelectedPosts(c *gin.Context) {
	logger := log.GetLogger()
	user := c.Param("user")
	var req []VelogPostRequestModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logger.Infof("download selected posts, user: %v", user)

	urlSlugList := []string{}
	for _, data := range req {
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

func getSeries(c *gin.Context) {
	user := c.Param("user")
	series, err := service.GetSeries(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": series,
	})
}

func getPostsInSeries(c *gin.Context) {
	logger := log.GetLogger()
	user := c.Param("user")
	seriesUrlSlug := c.Param("url_slug")
	logger.Infof("[readSeries] user: %v,seriesUrlSlug: %v", user, seriesUrlSlug)

	// TODO read series login
	postsInSeries, err := service.GetPostsInSereis(user, seriesUrlSlug)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}

	c.JSON(http.StatusOK, postsInSeries)
}

func downloadAllSeries(c *gin.Context) {
	c.String(http.StatusOK, "all series")
}

func downloadSelectedSeries(c *gin.Context) {
	c.String(http.StatusOK, "selected series")
}
