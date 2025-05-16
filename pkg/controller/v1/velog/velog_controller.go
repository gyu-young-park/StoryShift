package v1velogcontroller

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gyu-young-park/StoryShift/pkg/log"
	servicevelog "github.com/gyu-young-park/StoryShift/pkg/service/velog"
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
	router.Use(validateVelogUser())
	router.GET("/", checkVelogService)
	router.GET("/:user", getUserProfile)
	router.GET("/:user/post/:url_slug", getPost)
	router.GET("/:user/post/download", downloadPost)
	router.GET("/:user/posts", getPosts)
	router.POST("/:user/posts/download", downloadSelectedPosts)
	router.GET("/:user/posts/download", downloadAllPosts)
	router.GET("/:user/series", getSeries)
	router.GET("/:user/series/:url_slug", getPostsInSeries)
	router.GET("/:user/series/:url_slug/download", downloadSeries)
	router.GET("/:user/series/download", downloadAllSeries)
	router.POST("/:user/series/download", downloadSelectedSeries)
}

func checkVelogService(c *gin.Context) {
	if !servicevelog.IsVelogUserExists("velopert") {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "health"})
}

func getPost(c *gin.Context) {
	logger := log.GetLogger()
	user := c.Param("user")
	urlSlug := c.Param("url_slug")

	logger.Debugf("Velog username: %s, url_slug: %s",
		user,
		urlSlug)

	velogPost, err := servicevelog.GetPost(user, urlSlug)
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

	var req VelogGetPostsRequestModel
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Count == 0 {
		req.Count = 10
	}

	logger.Infof("Velog username: %s, post_id: %s, count: %v",
		user,
		req.PostId,
		req.Count)

	velogPosts, err := servicevelog.GetPosts(user, req.PostId, req.Count)
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

	var req VelogDownloadPostRequestModel
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Debugf("Velog username: %s, url_slug: %s", user, req.URLSlug)
	closeFunc, zipFilename, err := servicevelog.FetchVelogPostZip(user, req.URLSlug)
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

	closeFunc, zipFilename, err := servicevelog.FetchAllVelogPostsZip(user)
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
	var req []VelogDownloadPostRequestModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logger.Infof("download selected posts, user: %v", user)

	urlSlugList := []string{}
	for _, data := range req {
		urlSlugList = append(urlSlugList, data.URLSlug)
	}

	closeFunc, zipFilename, err := servicevelog.FetchSelectedVelogPostsZip(user, urlSlugList)
	defer closeFunc()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.FileAttachment(zipFilename, filepath.Base(zipFilename))
}

func getSeries(c *gin.Context) {
	user := c.Param("user")
	series, err := servicevelog.GetSeries(user)
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

	postsInSeries, err := servicevelog.GetPostsInSereis(user, seriesUrlSlug)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}

	c.JSON(http.StatusOK, postsInSeries)
}

func downloadSeries(c *gin.Context) {
	user := c.Param("user")
	seriesUrlSlug := c.Param("url_slug")

	closeFunc, zipFilename, err := servicevelog.FetchSeriesZip(user, seriesUrlSlug)
	defer closeFunc()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	c.FileAttachment(zipFilename, filepath.Base(zipFilename))
}

func downloadSelectedSeries(c *gin.Context) {
	user := c.Param("user")

	var req []VelogDownloadSelectedSeriesRequestModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}

	seriesURLSlugList := []string{}
	for _, seriesUrlSlug := range req {
		seriesURLSlugList = append(seriesURLSlugList, seriesUrlSlug.URLSlug)
	}

	closeFunc, zipFilename, err := servicevelog.FetchSelectedSeriesZip(user, seriesURLSlugList)
	defer closeFunc()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	c.FileAttachment(zipFilename, filepath.Base(zipFilename))
}

func downloadAllSeries(c *gin.Context) {
	user := c.Param("user")

	closeFunc, zipFilename, err := servicevelog.FetchAllSeriesZip(user)
	defer closeFunc()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	c.FileAttachment(zipFilename, filepath.Base(zipFilename))
}

func getUserProfile(c *gin.Context) {
	user := c.Param("user")

	userProfile, err := servicevelog.GetUserProfile(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"user_profile": userProfile,
	})
}
