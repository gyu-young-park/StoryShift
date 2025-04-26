package v1controller

type VelogUserNameReqestModel struct {
	Name string `form:"name" binding:"required"`
}

type VelogPostRequestModel struct {
	VelogUserNameReqestModel
	URLSlog string `form:"url_slog"`
}

type VelogPostsRequestModel struct {
	VelogUserNameReqestModel
	PostId string `form:"post_id"`
	Limit  int    `form:"limit"`
}
