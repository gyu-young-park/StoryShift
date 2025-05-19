package velog

import (
	"fmt"
	"time"
)

type commonVelogPost struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *commonVelogPost) mapped(id, title string, releasedAt, updatedAt time.Time) error {
	if id == "" {
		return ErrNoMatchPost
	}
	c.ID = id
	c.Title = title
	c.CreatedAt = releasedAt
	c.UpdatedAt = updatedAt
	return nil
}

type VelogPost struct {
	commonVelogPost
	Body string `json:"body"`
}

func (p *VelogPost) mapped(model readPostModel) error {
	err := p.commonVelogPost.mapped(model.Data.Post.ID, model.Data.Post.Title,
		model.Data.Post.ReleasedAt, model.Data.Post.UpdatedAt,
	)

	if err != nil {
		return err
	}

	p.Body = model.Data.Post.Body
	return nil
}

type VelogPostsItemList []VelogPostsItem

func (pl *VelogPostsItemList) mapped(model postsModel) error {
	for i, _ := range model.Data.Posts {
		postsItem := VelogPostsItem{}
		err := postsItem.mapped(model, i)
		if err != nil {
			return err
		}
		*pl = append(*pl, postsItem)
	}

	fmt.Println(pl)

	return nil
}

type VelogPostsItem struct {
	commonVelogPost
	ShortDesc string   `json:"short_description"`
	Thumnail  string   `json:"thumnail"`
	UrlSlug   string   `json:"url_slug"`
	Tags      []string `json:"tags"`
}

func (p *VelogPostsItem) mapped(model postsModel, index int) error {
	post := model.Data.Posts[index]
	err := p.commonVelogPost.mapped(post.ID, post.Title,
		post.ReleasedAt, post.UpdatedAt,
	)

	if err != nil {
		return err
	}

	p.ShortDesc = post.ShortDescription
	p.Thumnail = post.Thumbnail
	p.UrlSlug = post.URLSlug
	p.Tags = post.Tags

	return nil
}

type VelogSeriesBase struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (sb *VelogSeriesBase) mapped(id, name string) error {
	if id == "" {
		return ErrNoMatchPost
	}

	sb.ID = id
	sb.Name = name

	return nil
}

type VelogSeriesItemList []VelogSeriesItem

func (sl *VelogSeriesItemList) mapped(model userSeriesListModel) error {
	for i, _ := range model.Data.User.SeriesList {
		seriesItem := VelogSeriesItem{}
		err := seriesItem.mapped(model, i)
		if err != nil {
			return err
		}
		*sl = append(*sl, seriesItem)
	}

	return nil
}

type VelogSeriesItem struct {
	VelogSeriesBase
	URLSlug   string    `json:"url_slug"`
	Count     int       `json:"count"`
	Thumbnail string    `json:"thumbnail"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (si *VelogSeriesItem) mapped(model userSeriesListModel, index int) error {
	series := model.Data.User.SeriesList[index]
	err := si.VelogSeriesBase.mapped(series.ID, series.Name)
	if err != nil {
		return err
	}

	si.URLSlug = series.URLSlug
	si.Count = series.PostsCount
	si.Thumbnail = series.Thumbnail
	si.UpdatedAt = series.UpdatedAt

	return nil
}

type velogReadSeriesItem struct {
	Title     string    `json:"title"`
	URLSlug   string    `json:"url_slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (rsi *velogReadSeriesItem) mapped(model readSeriesModel, index int) error {
	seriesPost := model.Data.Series.SeriesPosts[index]
	rsi.Title = seriesPost.Post.Title
	rsi.URLSlug = seriesPost.Post.URLSlug
	rsi.CreatedAt = seriesPost.Post.ReleasedAt
	rsi.UpdatedAt = seriesPost.Post.UpdatedAt

	return nil
}

type VelogReadSeries struct {
	VelogSeriesBase
	Posts []velogReadSeriesItem `json:"posts"`
}

func (rs *VelogReadSeries) mapped(model readSeriesModel) error {
	err := rs.VelogSeriesBase.mapped(model.Data.Series.ID, model.Data.Series.Name)
	if err != nil {
		return err
	}

	for i, _ := range model.Data.Series.SeriesPosts {
		seriesItem := velogReadSeriesItem{}
		err := seriesItem.mapped(model, i)
		if err != nil {
			return err
		}
	}

	return nil
}

type VelogUserProfile struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Describe  string `json:"describe"`
	Thumbnail string `json:"thumbnail"`
	Bio       string `json:"Bio"`
}

func (up *VelogUserProfile) mapped(model userProfileModel) error {
	if model.Data.User.ID == "" {
		return ErrNoMatchUser
	}

	up.Id = model.Data.User.ID
	up.Username = model.Data.User.Username
	up.Describe = model.Data.User.Profile.DisplayName
	up.Thumbnail = model.Data.User.Profile.Thumbnail
	up.Bio = model.Data.User.Profile.ShortBio

	return nil
}
