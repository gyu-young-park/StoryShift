package velog

import "time"

type commonVelogPost struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VelogPost struct {
	commonVelogPost
	Body string `json:"body"`
}

type VelogPostsItem struct {
	commonVelogPost
	ShortDesc string   `json:"short_description"`
	Thumnail  string   `json:"thumnail"`
	UrlSlug   string   `json:"url_slug"`
	Tags      []string `json:"tags"`
}

type VelogSeriesBase struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type VelogSeriesItem struct {
	VelogSeriesBase
	URLSlug   string    `json:"url_slug"`
	Count     int       `json:"count"`
	Thumbnail string    `json:"thumbnail"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VelogReadSeriesItem struct {
	Title     string    `json:"title"`
	URLSlug   string    `json:"url_slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VelogReadSeries struct {
	VelogSeriesBase
	Posts []VelogReadSeriesItem `json:"posts"`
}
