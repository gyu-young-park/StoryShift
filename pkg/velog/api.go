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

func (v VelogAPI) GetPost(urlSlug string) (VelogPost, error) {
	reqBody := graphQLQuery.readPost(v.Username, urlSlug)

	resp, err := httpclient.Post(httpclient.PostRequestParam{
		URL:         v.VelogAPIURL,
		Body:        bytes.NewBuffer([]byte(reqBody)),
		ContentType: "application/json",
	})

	var model readPostModel
	if err != nil {
		return VelogPost{}, err
	}

	err = json.Unmarshal(resp, &model)
	if err != nil {
		return VelogPost{}, err
	}

	post := VelogPost{
		ID:        model.Data.Post.ID,
		Title:     model.Data.Post.Title,
		Body:      model.Data.Post.Body,
		CreatedAt: model.Data.Post.ReleasedAt,
		UpdatedAt: model.Data.Post.UpdatedAt,
	}

	return post, nil
}

// func (v *VelogAPI) GetAllPost(cursor string) []VelogPost {

// }
