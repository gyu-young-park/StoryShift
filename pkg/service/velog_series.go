package service

import (
	"os"

	"github.com/gyu-young-park/StoryShift/internal/config"
	"github.com/gyu-young-park/StoryShift/pkg/file"
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

func FetchSeriesZip(username string, seriesUrlSlug string) (closeFunc, string, error) {
	fileHandler := file.NewFileHandler()
	closeFunc := func() {
		defer fileHandler.Close()
	}

	postsInSeriesModel, err := GetPostsInSereis(username, seriesUrlSlug)
	if err != nil {
		return closeFunc, "", err
	}

	fileList := []*os.File{}
	for _, post := range postsInSeriesModel.Posts {
		postFilePath, err := fileHandler.CreateFile(file.File{
			FileMeta: file.FileMeta{
				Name:      post.Title,
				Extention: "md",
			},
			Content: post.Body,
		})

		if err != nil {
			return closeFunc, "", err
		}
		fileList = append(fileList, fileHandler.GetFileWithLocked(postFilePath))
	}

	zipFileName, err := fileHandler.CreateZipFile(file.ZipFile{
		FileMeta: file.FileMeta{
			Name:      postsInSeriesModel.Name,
			Extention: "zip",
		},
		Files: fileList,
	})

	if err != nil {
		return closeFunc, "", err
	}

	return closeFunc, zipFileName, nil

}

func FetchSelectedSeriesZip(username string, seriesUrlSlugList []string) (closeFunc, string, error) {
	// TODO
	return func() {}, "", nil
}

func FetchAllSeriesZip(username string) {
	// TODO
}
