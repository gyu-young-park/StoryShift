package v1controller

type VelogUserNameReqestModel struct {
	Name string `form:"name" json:"name" binding:"required"`
}

type VelogPostRequestModel struct {
	VelogUserNameReqestModel
	URLSlug string `form:"url_slug" json:"url_slug"`
}

type VelogPostsRequestModel struct {
	VelogUserNameReqestModel
	PostId string `form:"post_id" json:"post_id"`
	Count  int    `form:"count" json:"count"`
}
