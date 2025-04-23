package velog

import (
	"bytes"
	"encoding/json"

	"github.com/gyu-young-park/VelogStoryShift/internal/httpclient"
)

type VelogAPI struct {
	VelogAPIURL string
	Username    string
}

func NewVelogAPI(apiUrl string, username string) VelogAPI {
	return VelogAPI{
		VelogAPIURL: apiUrl,
		Username:    username,
	}
}

func (v VelogAPI) GetPost(postId string) (VelogAPIGetPostRespModel, error) {
	reqBody := getLastPostHistoryRawGraphQL(postId) // "4023bf7e-df1c-4288-9e4f-a37983406912"

	resp, err := httpclient.Post(httpclient.PostRequestParam{
		URL:         v.VelogAPIURL, // "https://v2.velog.io/graphql"
		Body:        bytes.NewBuffer([]byte(reqBody)),
		ContentType: "application/json",
	})

	var velogAPIGetPostRespModel VelogAPIGetPostRespModel
	if err != nil {
		return velogAPIGetPostRespModel, err
	}

	err = json.Unmarshal(resp, &velogAPIGetPostRespModel)
	if err != nil {
		return velogAPIGetPostRespModel, err
	}

	return velogAPIGetPostRespModel, nil
}
