package v1controller

type VelogPostRequestModel struct {
	URLSlug string `form:"url_slug" json:"url_slug"`
}

type VelogPostsRequestModel struct {
	PostId string `form:"post_id" json:"post_id"`
	Count  int    `form:"count" json:"count"`
}

type VelogReadSeriesRequestModel struct {
	// postid -> postid부터 ~ count까지 가져오도록 하기
	Count int `form:"count" json:"count"`
}
