package service

import (
	"fmt"
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

func fetchSeries(fileHandler *file.FileHandler, username string, seriesUrlSlug string) ([]*os.File, error) {
	postsInSeriesModel, err := GetPostsInSereis(username, seriesUrlSlug)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		fileList = append(fileList, fileHandler.GetFileWithLocked(postFilePath))
	}

	return fileList, nil
}

func FetchSeriesZip(username string, seriesUrlSlug string) (closeFunc, string, error) {
	fileHandler := file.NewFileHandler()
	closeFunc := func() {
		defer fileHandler.Close()
	}

	zipfileList, err := fetchSeries(fileHandler, username, seriesUrlSlug)
	if err != nil {
		return closeFunc, "", err
	}

	zipFileName, err := fileHandler.CreateZipFile(file.ZipFile{
		FileMeta: file.FileMeta{
			Name:      fmt.Sprintf("%s-%s", username, seriesUrlSlug),
			Extention: "zip",
		},
		Files: zipfileList,
	})

	if err != nil {
		return closeFunc, "", err
	}

	return closeFunc, zipFileName, nil

}

func FetchSelectedSeriesZip(username string, seriesUrlSlugList []string) (closeFunc, string, error) {
	fileHandler := file.NewFileHandler()
	closeFunc := func() {
		defer fileHandler.Close()
	}

	zipfileList := []*os.File{}
	for _, seriesUrlSlug := range seriesUrlSlugList {
		fileList, err := fetchSeries(fileHandler, username, seriesUrlSlug)
		if err != nil {
			return closeFunc, "", err
		}

		// refactor: series data 가져오기, 공통 로직 분리하기
		zipFileName, err := fileHandler.CreateZipFile(file.ZipFile{
			FileMeta: file.FileMeta{
				Name:      seriesUrlSlug,
				Extention: "zip",
			},
			Files: fileList,
		})

		if err != nil {
			return closeFunc, "", err
		}

		zipFh := fileHandler.GetFileWithLocked(zipFileName)
		zipFh.Seek(0, 0)

		zipfileList = append(zipfileList, zipFh)
	}

	zipFileName, err := fileHandler.CreateZipFile(file.ZipFile{
		FileMeta: file.FileMeta{
			Name:      fmt.Sprintf("%s-%s", username, "series-post-list"),
			Extention: "zip",
		},
		Files: zipfileList,
	})

	if err != nil {
		return closeFunc, "", err
	}

	return closeFunc, zipFileName, nil
}

func FetchAllSeriesZip(username string) (closeFunc, string, error) {
	seriesItemList, err := GetSeries(username)
	if err != nil {
		return func() {}, "", err
	}

	fileHandler := file.NewFileHandler()
	closeFunc := func() {
		defer fileHandler.Close()
	}

	zipfileList := []*os.File{}
	for _, seiriesItem := range seriesItemList {
		fileList, err := fetchSeries(fileHandler, username, seiriesItem.URLSlug)
		if err != nil {
			return closeFunc, "", err
		}

		zipFileName, err := fileHandler.CreateZipFile(file.ZipFile{
			FileMeta: file.FileMeta{
				Name:      seiriesItem.URLSlug,
				Extention: "zip",
			},
			Files: fileList,
		})

		if err != nil {
			return closeFunc, "", err
		}
		zipFh := fileHandler.GetFileWithLocked(zipFileName)
		zipFh.Seek(0, 0)

		zipfileList = append(zipfileList, zipFh)
	}

	zipFileName, err := fileHandler.CreateZipFile(file.ZipFile{
		FileMeta: file.FileMeta{
			Name:      username + "-all-series",
			Extention: "zip",
		},
		Files: zipfileList,
	})

	if err != nil {
		return closeFunc, "", err
	}

	return closeFunc, zipFileName, err
}
