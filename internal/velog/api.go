package velog

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/gyu-young-park/VelogStoryShift/internal/httpclient"
)

type VelogAPI struct {
}

type variables struct {
	Username string `json:"username"`
}

type GraphQLReqBody struct {
	OperationName string    `json:"operationName"`
	Variables     variables `json:"variables"`
	Query         string    `json:"query"`
}

func GetList() {
	reqBody := GraphQLReqBody{
		OperationName: "UserSeriesList",
		Variables: variables{
			Username: "chappi",
		},
		Query: "query UserSeriesList($username: String!) {\n  user(username: $username) {\n    id\n    series_list {\n      id\n      name\n      description\n      url_slug\n      thumbnail\n      updated_at\n      posts_count\n      __typename\n    }\n    __typename\n  }\n}\n",
	}
	buf, _ := json.Marshal(reqBody)
	resp, err := httpclient.Post("https://v2.velog.io/graphql", bytes.NewBuffer(buf))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}

func GetPostData() {

}
