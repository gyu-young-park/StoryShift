package v1controller

type VelogUserNameReqestModel struct {
	Name string `form:"name" binding:"required"`
}

type VelogPostRequestModel struct {
	VelogUserNameReqestModel
	URLSlug string `form:"url_slug"`
}

type VelogPostsRequestModel struct {
	VelogUserNameReqestModel
	PostId string `form:"post_id"`
	Count  int    `form:"count"`
}
