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

func (v VelogAPI) Posts(cursor string, limit int) ([]VelogPostsItem, error) {
	reqBody := graphQLQuery.posts(v.Username, cursor, limit)

	resp, err := httpclient.Post(httpclient.PostRequestParam{
		URL:         v.VelogAPIURL,
		Body:        bytes.NewBuffer([]byte(reqBody)),
		ContentType: "application/json",
	})

	if err != nil {
		return []VelogPostsItem{}, err
	}

	var model postsModel
	err = json.Unmarshal(resp, &model)
	if err != nil {
		return []VelogPostsItem{}, err
	}

	posts := []VelogPostsItem{}

	for _, post := range model.Data.Posts {
		posts = append(posts, VelogPostsItem{
			commonVelogPost: commonVelogPost{
				ID:        post.ID,
				Title:     post.Title,
				CreatedAt: post.ReleasedAt,
				UpdatedAt: post.UpdatedAt,
			},
			ShortDesc: post.ShortDescription,
			Thumnail:  post.Thumbnail,
			UrlSlog:   post.URLSlug,
			Tags:      post.Tags,
		})
	}

	return posts, nil
}

func (v VelogAPI) Post(urlSlug string) (VelogPost, error) {
	reqBody := graphQLQuery.readPost(v.Username, urlSlug)

	resp, err := httpclient.Post(httpclient.PostRequestParam{
		URL:         v.VelogAPIURL,
		Body:        bytes.NewBuffer([]byte(reqBody)),
		ContentType: "application/json",
	})

	if err != nil {
		return VelogPost{}, err
	}

	var model readPostModel
	err = json.Unmarshal(resp, &model)
	if err != nil {
		return VelogPost{}, err
	}

	post := VelogPost{
		commonVelogPost: commonVelogPost{
			ID:        model.Data.Post.ID,
			Title:     model.Data.Post.Title,
			CreatedAt: model.Data.Post.ReleasedAt,
			UpdatedAt: model.Data.Post.UpdatedAt,
		},
		Body: model.Data.Post.Body,
	}

	return post, nil
}
