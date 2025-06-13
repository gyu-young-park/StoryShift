package v1velogcontroller

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gyu-young-park/StoryShift/pkg/log"
	servicevelog "github.com/gyu-young-park/StoryShift/pkg/service/velog"
)

type VelogController struct {
	service  *servicevelog.VelogService
	APIGroup string
}

func NewVelogController(service *servicevelog.VelogService) *VelogController {
	return &VelogController{
		service:  service,
		APIGroup: "/velog",
	}
}

func (v *VelogController) GetAPIGroup() string {
	return v.APIGroup
}

func (v *VelogController) RegisterAPI(router *gin.RouterGroup) {
	router.Use(validateVelogUser(v.service))
	router.GET("/", v.checkVelogService)
	router.GET("/:user", v.getUserProfile)
	router.GET("/:user/post/:url_slug", v.getPost)
	router.GET("/:user/post/download", v.downloadPost)
	router.GET("/:user/posts", v.getPosts)
	router.POST("/:user/posts/download", v.downloadSelectedPosts)
	router.GET("/:user/posts/download", v.downloadAllPosts)
	router.GET("/:user/series", v.getSeries)
	router.GET("/:user/series/:url_slug", v.getPostsInSeries)
	router.GET("/:user/series/:url_slug/download", v.downloadSeries)
	router.GET("/:user/series/download", v.downloadAllSeries)
	router.POST("/:user/series/download", v.downloadSelectedSeries)
}

func (v *VelogController) checkVelogService(c *gin.Context) {
	if !v.service.IsVelogUserExists("velopert") {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "health"})
}

func (v *VelogController) getPost(c *gin.Context) {
	logger := log.GetLogger()
	user := c.Param("user")
	urlSlug := c.Param("url_slug")

	logger.Debugf("Velog username: %s, url_slug: %s",
		user,
		urlSlug)

	velogPost, err := v.service.GetPost(user, urlSlug)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": velogPost,
	})
}

func (v *VelogController) getPosts(c *gin.Context) {
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

	velogPosts, err := v.service.GetPosts(user, req.PostId, req.Count)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": velogPosts,
	})
}

func (v *VelogController) downloadPost(c *gin.Context) {
	logger := log.GetLogger()
	user := c.Param("user")

	var req VelogDownloadPostRequestModel
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Debugf("Velog username: %s, url_slug: %s", user, req.URLSlug)
	closeFunc, zipFilename, err := v.service.FetchVelogPostZip(user, req.URLSlug)
	defer closeFunc()

	logger.Infof("get zip file: %s", zipFilename)
	if err != nil {
		logger.Errorf("failed to return zip file: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.FileAttachment(zipFilename, filepath.Base(zipFilename))
}

func (v *VelogController) downloadAllPosts(c *gin.Context) {
	logger := log.GetLogger()
	user := c.Param("user")

	var req VelogDownloadAllPostRequestModel
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	closeFunc, zipFilename, err := v.service.FetchAllVelogPostsZip(user, req.Refresh, req.Image)
	defer closeFunc()
	if err != nil {
		logger.Errorf("server error occureed: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.FileAttachment(zipFilename, filepath.Base(zipFilename))
}

func (v *VelogController) downloadSelectedPosts(c *gin.Context) {
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

	closeFunc, zipFilename, err := v.service.FetchSelectedVelogPostsZip(user, urlSlugList)
	defer closeFunc()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.FileAttachment(zipFilename, filepath.Base(zipFilename))
}

func (v *VelogController) getSeries(c *gin.Context) {
	user := c.Param("user")
	series, err := v.service.GetSeries(user, false)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": series,
	})
}

func (v *VelogController) getPostsInSeries(c *gin.Context) {
	logger := log.GetLogger()
	user := c.Param("user")
	seriesUrlSlug := c.Param("url_slug")
	logger.Infof("[readSeries] user: %v,seriesUrlSlug: %v", user, seriesUrlSlug)

	postsInSeries, err := v.service.GetPostsInSereis(user, seriesUrlSlug, false)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}

	c.JSON(http.StatusOK, postsInSeries)
}

func (v *VelogController) downloadSeries(c *gin.Context) {
	user := c.Param("user")
	seriesUrlSlug := c.Param("url_slug")

	closeFunc, zipFilename, err := v.service.FetchSeriesZip(user, seriesUrlSlug)
	defer closeFunc()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	c.FileAttachment(zipFilename, filepath.Base(zipFilename))
}

func (v *VelogController) downloadSelectedSeries(c *gin.Context) {
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

	closeFunc, zipFilename, err := v.service.FetchSelectedSeriesZip(user, seriesURLSlugList)
	defer closeFunc()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	c.FileAttachment(zipFilename, filepath.Base(zipFilename))
}

func (v *VelogController) downloadAllSeries(c *gin.Context) {
	user := c.Param("user")
	refresh := c.Query("refresh")
	isRefresh, err := strconv.ParseBool(refresh)
	if err != nil {
		isRefresh = false
	}

	closeFunc, zipFilename, err := v.service.FetchAllSeriesZip(user, isRefresh)
	defer closeFunc()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	c.FileAttachment(zipFilename, filepath.Base(zipFilename))
}

func (v *VelogController) getUserProfile(c *gin.Context) {
	user := c.Param("user")

	userProfile, err := v.service.GetUserProfile(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"user_profile": userProfile,
	})
}
