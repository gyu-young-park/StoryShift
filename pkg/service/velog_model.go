package service

import "github.com/gyu-young-park/StoryShift/pkg/velog"

type PostsInSeriesModel struct {
	velog.VelogSeriesBase
	Posts []velog.VelogPost `json:"posts"`
}
