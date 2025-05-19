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

func (v VelogAPI) Posts(cursor string, limit int) (VelogPostsItemList, error) {
	reqBody := graphQLQuery.posts(v.Username, cursor, limit)

	resp, err := httpclient.Post(httpclient.PostRequestParam{
		URL:         v.VelogAPIURL,
		Body:        bytes.NewBuffer([]byte(reqBody)),
		ContentType: "application/json",
	})

	if err != nil {
		return VelogPostsItemList{}, err
	}

	var model postsModel
	err = json.Unmarshal(resp.Body, &model)
	if err != nil {
		return VelogPostsItemList{}, err
	}

	posts := VelogPostsItemList{}
	err = posts.mapped(model)
	return posts, err
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
	err = json.Unmarshal(resp.Body, &model)
	if err != nil {
		return VelogPost{}, err
	}

	post := VelogPost{}
	err = post.mapped(model)
	if err != nil {
		return post, err
	}

	return post, nil
}

func (v VelogAPI) Series() (VelogSeriesItemList, error) {
	reqBody := graphQLQuery.userSeriesList(v.Username)
	resp, err := httpclient.Post(httpclient.PostRequestParam{
		URL:         v.VelogAPIURL,
		Body:        bytes.NewBuffer([]byte(reqBody)),
		ContentType: "application/json",
	})

	if err != nil {
		return VelogSeriesItemList{}, err
	}

	var model userSeriesListModel
	err = json.Unmarshal(resp.Body, &model)
	if err != nil {
		return VelogSeriesItemList{}, err
	}

	seriesList := VelogSeriesItemList{}
	err = seriesList.mapped(model)
	return seriesList, err
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
	err = json.Unmarshal(resp.Body, &model)
	if err != nil {
		return VelogReadSeries{}, err
	}

	readSeries := VelogReadSeries{}
	err = readSeries.mapped(model)
	return readSeries, err
}

func (v VelogAPI) UserProfile() (VelogUserProfile, error) {
	reqBody := graphQLQuery.userProfile(v.Username)
	resp, err := httpclient.Post(httpclient.PostRequestParam{
		URL:         v.VelogAPIURL,
		Body:        bytes.NewBuffer([]byte(reqBody)),
		ContentType: "application/json",
	})

	if err != nil {
		return VelogUserProfile{}, err
	}

	var model userProfileModel
	err = json.Unmarshal(resp.Body, &model)

	if err != nil {
		return VelogUserProfile{}, err
	}

	velogUserProfile := VelogUserProfile{}
	err = velogUserProfile.mapped(model)
	return velogUserProfile, err
}
