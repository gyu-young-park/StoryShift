package service

import (
	"github.com/gyu-young-park/StoryShift/internal/config"
	"github.com/gyu-young-park/StoryShift/pkg/velog"
)

func GetSeries(username string) ([]velog.VelogSeriesItem, error) {
	velogApi := velog.NewVelogAPI(config.Manager.VelogConfig.URL, username)

	seriesList, err := velogApi.Series()
	if err != nil {
		return []velog.VelogSeriesItem{}, err
	}

	return seriesList, nil
}

func GetPostsInSereis(username, seriesUrlSlug string) (PostsInSeriesModel, error) {
	velogApi := velog.NewVelogAPI(config.Manager.VelogConfig.URL, username)
	readSeriesList, err := velogApi.ReadSeries(seriesUrlSlug)
	if err != nil {
		return PostsInSeriesModel{}, err
	}

	postInSeriesModel := PostsInSeriesModel{
		VelogSeriesBase: readSeriesList.VelogSeriesBase,
		Posts:           []velog.VelogPost{},
	}
	for _, postInSeries := range readSeriesList.Posts {
		post, err := velogApi.Post(postInSeries.URLSlug)
		if err != nil {
			return PostsInSeriesModel{}, nil
		}
		postInSeriesModel.Posts = append(postInSeriesModel.Posts, post)
	}

	return postInSeriesModel, nil
}

func DownloadSelectedSeries(username string, seriesId []string) {
	// TODO
}

func DownloadAllSeries(username string) {
	// TODO
}
