package v1controller

import (
	"github.com/gin-gonic/gin"
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
	//
}

func getPosts(c *gin.Context) {
}

func getPost(c *gin.Context) {

}
