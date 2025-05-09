package v1controller

type VelogPostRequestModel struct {
	URLSlug string `form:"url_slug" json:"url_slug"`
}

type VelogPostsRequestModel struct {
	PostId string `form:"post_id" json:"post_id"`
	Count  int    `form:"count" json:"count"`
}
