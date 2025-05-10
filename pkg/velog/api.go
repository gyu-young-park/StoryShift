package velog

import (
	"bytes"
	"encoding/json"

	"github.com/gyu-young-park/StoryShift/internal/httpclient"
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
			UrlSlug:   post.URLSlug,
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

func (v VelogAPI) Series() ([]VelogSeriesItem, error) {
	reqBody := graphQLQuery.userSeriesList(v.Username)
	resp, err := httpclient.Post(httpclient.PostRequestParam{
		URL:         v.VelogAPIURL,
		Body:        bytes.NewBuffer([]byte(reqBody)),
		ContentType: "application/json",
	})

	if err != nil {
		return []VelogSeriesItem{}, err
	}

	var model userSeriesListModel
	err = json.Unmarshal(resp, &model)
	if err != nil {
		return []VelogSeriesItem{}, err
	}

	seriesList := []VelogSeriesItem{}
	for _, series := range model.Data.User.SeriesList {
		seriesList = append(seriesList, VelogSeriesItem{
			VelogSeriesBase: VelogSeriesBase{
				ID:   series.ID,
				Name: series.Name,
			},
			URLSlug:   series.URLSlug,
			Count:     series.PostsCount,
			Thumbnail: series.Thumbnail,
			UpdatedAt: series.UpdatedAt,
		})
	}

	return seriesList, nil
}

func (v VelogAPI) ReadSeries(urlSlug string) (VelogReadSeries, error) {
	reqBody := graphQLQuery.readSeries(v.Username, urlSlug)
	resp, err := httpclient.Post(httpclient.PostRequestParam{
		URL:         v.VelogAPIURL,
		Body:        bytes.NewBuffer([]byte(reqBody)),
		ContentType: "application/json",
	})

	if err != nil {
		return VelogReadSeries{}, err
	}

	var model readSeriesModel
	err = json.Unmarshal(resp, &model)
	if err != nil {
		return VelogReadSeries{}, err
	}

	readSeries := VelogReadSeries{
		VelogSeriesBase: VelogSeriesBase{
			ID:   model.Data.Series.ID,
			Name: model.Data.Series.Name,
		},
		Posts: []velogReadSeriesItem{},
	}
	for _, post := range model.Data.Series.SeriesPosts {
		readSeries.Posts = append(readSeries.Posts, velogReadSeriesItem{
			Title:     post.Post.Title,
			URLSlug:   post.Post.URLSlug,
			CreatedAt: post.Post.ReleasedAt,
			UpdatedAt: post.Post.UpdatedAt,
		})
	}

	return readSeries, nil
}
