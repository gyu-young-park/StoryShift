package servicevelog

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gyu-young-park/StoryShift/internal/cache"
	"github.com/gyu-young-park/StoryShift/pkg/file"
	"github.com/gyu-young-park/StoryShift/pkg/log"
	"github.com/gyu-young-park/StoryShift/pkg/velog"
	"github.com/gyu-young-park/StoryShift/pkg/worker"
)

func (v *VelogService) GetSeries(username string, isRefresh bool) (velog.VelogSeriesItemList, error) {
	seriesList, err := v.cacheManager.CallWithCache(cache.CacheOptBuilder.
		Timeout(time.Second*2).
		TTL(time.Minute*10).
		Refresh(isRefresh).
		Build(fmt.Sprintf("%s-%s", username, "series")),
		func() (string, error) {
			series, err := v.velogAPI.Series(username)
			if err != nil {
				return "", err
			}

			bSeries, _ := json.Marshal(series)
			return string(bSeries), nil
		})

	var seriesListModel velog.VelogSeriesItemList
	if err != nil {
		return seriesListModel, err
	}

	json.Unmarshal([]byte(seriesList), &seriesListModel)
	return seriesListModel, nil
}

func (v *VelogService) GetPostsInSereis(username, seriesUrlSlug string, isRefresh bool) (PostsInSeriesModel, error) {
	postInSeries, err := v.cacheManager.CallWithCache(cache.CacheOptBuilder.
		Timeout(time.Second*2).
		TTL(time.Minute*10).
		Refresh(isRefresh).
		Build(fmt.Sprintf("%s-%s-%s", "series", username, seriesUrlSlug)),
		func() (string, error) {
			readSeriesList, err := v.velogAPI.ReadSeries(username, seriesUrlSlug)
			if err != nil {
				return "", err
			}

			postInSeriesModel := PostsInSeriesModel{
				VelogSeriesBase: readSeriesList.VelogSeriesBase,
				Posts:           []velog.VelogPost{},
			}

			for _, postInSeries := range readSeriesList.Posts {
				post, err := v.velogAPI.Post(username, postInSeries.URLSlug)
				if err != nil {
					return "", nil
				}
				postInSeriesModel.Posts = append(postInSeriesModel.Posts, post)
			}

			bPostInSeriesModel, _ := json.Marshal(postInSeriesModel)
			return string(bPostInSeriesModel), nil
		})

	var postInSeriesModel PostsInSeriesModel
	if err != nil {
		return postInSeriesModel, err
	}

	json.Unmarshal([]byte(postInSeries), &postInSeriesModel)
	return postInSeriesModel, nil
}

func (v *VelogService) fetchSeries(fileHandler *file.FileHandler, username string, seriesUrlSlug string, isRefresh bool) ([]*os.File, error) {
	postsInSeriesModel, err := v.GetPostsInSereis(username, seriesUrlSlug, isRefresh)
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

func (v *VelogService) FetchSeriesZip(username string, seriesUrlSlug string) (closeFunc, string, error) {
	fileHandler := file.NewFileHandler()
	closeFunc := func() {
		defer fileHandler.Close()
	}

	zipfileList, err := v.fetchSeries(fileHandler, username, seriesUrlSlug, false)
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

func (v *VelogService) FetchSelectedSeriesZip(username string, seriesUrlSlugList []string) (closeFunc, string, error) {
	fileHandler := file.NewFileHandler()
	closeFunc := func() {
		defer fileHandler.Close()
	}

	zipfileList := []*os.File{}
	for _, seriesUrlSlug := range seriesUrlSlugList {
		fileList, err := v.fetchSeries(fileHandler, username, seriesUrlSlug, false)
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

func (v *VelogService) FetchAllSeriesZip(username string, isRefresh bool) (closeFunc, string, error) {
	logger := log.GetLogger()
	seriesItemList, err := v.GetSeries(username, isRefresh)
	if err != nil {
		return func() {}, "", err
	}

	fileHandler := file.NewFileHandler()
	closeFunc := func() {
		defer fileHandler.Close()
	}

	ctx, cancel := context.WithCancel(context.Background())
	workerManager := worker.NewWorkerManager[velog.VelogSeriesItem, *os.File](ctx, fmt.Sprintf("%s-%s", "velog-series-zip", username), 5)
	defer cancel()
	defer workerManager.Close()

	zfhList := workerManager.Aggregate(seriesItemList, func(vsi velog.VelogSeriesItem) *os.File {
		fileList, err := v.fetchSeries(fileHandler, username, vsi.URLSlug, isRefresh)
		if err != nil {
			logger.Errorf("failed to fetch sereis: %v", err.Error())
			return nil
		}

		zipFileName, err := fileHandler.CreateZipFile(file.ZipFile{
			FileMeta: file.FileMeta{
				Name:      vsi.URLSlug,
				Extention: "zip",
			},
			Files: fileList,
		})

		if err != nil {
			logger.Errorf("failed to create sereis zip file: %v", err.Error())
			return nil
		}

		zipFh := fileHandler.GetFileWithLocked(zipFileName)
		zipFh.Seek(0, 0)

		return zipFh
	})

	zipfileList := []*os.File{}
	for _, zfh := range zfhList {
		if zfh != nil {
			zipfileList = append(zipfileList, zfh)
		}
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
