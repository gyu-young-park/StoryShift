package velog

import "time"

type variables struct {
	PostId   string `json:"post_id"`
	Username string `json:"username"`
}

type graphQLRequestBody struct {
	OperationName string    `json:"operationName"`
	Variables     variables `json:"variables"`
	Query         string    `json:"query"`
}

type VelogAPIGetPostRespModel struct {
	Data struct {
		LastPostHistory struct {
			ID         string    `json:"id"`
			Title      string    `json:"title"`
			Body       string    `json:"body"`
			CreatedAt  time.Time `json:"created_at"`
			IsMarkdown bool      `json:"is_markdown"`
			Typename   string    `json:"__typename"`
		} `json:"lastPostHistory"`
	} `json:"data"`
}
