package servicevelog

import (
	"context"
	"fmt"
	"os"

	"github.com/gyu-young-park/StoryShift/internal/config"
	"github.com/gyu-young-park/StoryShift/pkg/file"
	"github.com/gyu-young-park/StoryShift/pkg/log"
	"github.com/gyu-young-park/StoryShift/pkg/velog"
	"github.com/gyu-young-park/StoryShift/pkg/worker"
)

func GetPost(username, urlSlug string) (velog.VelogPost, error) {
	velogApi := velog.NewVelogAPI(config.Manager.VelogConfig.ApiUrl, username)
	velogPost, err := velogApi.Post(urlSlug)
	if err != nil {
		return velog.VelogPost{}, err
	}

	return velogPost, nil
}

func GetPosts(username, postId string, count int) (velog.VelogPostsItemList, error) {
	velogApi := velog.NewVelogAPI(config.Manager.VelogConfig.ApiUrl, username)
	posts, err := velogApi.Posts(postId, count)
	if err != nil {
		return velog.VelogPostsItemList{}, err
	}

	return posts, nil
}

type closeFunc func()

func FetchVelogPostZip(username, postId string) (closeFunc, string, error) {
	velogApi := velog.NewVelogAPI(config.Manager.VelogConfig.ApiUrl, username)
	velogPost, err := velogApi.Post(postId)

	if err != nil {
		return func() {}, "", err
	}

	f := file.File{
		FileMeta: file.FileMeta{
			Name:      velogPost.Title,
			Extention: "md",
		},
		Content: velogPost.Body,
	}

	fileHandler := file.NewFileHandler()
	zipFilename, err := fileHandler.MakeZipFileWithTitle(f)
	if err != nil {
		return func() {}, "", err
	}

	return func() {
		fileHandler.Close()
	}, zipFilename, err
}

type RenamedFileJSON struct {
	Origin string `json:"origin"`
	Rename string `json:"rename"`
}

func FetchAllVelogPostsZip(username string) (closeFunc, string, error) {
	logger := log.GetLogger()
	velogApi := velog.NewVelogAPI(config.Manager.VelogConfig.ApiUrl, username)
	fileHandler := file.NewFileHandler()

	closeFunc := func() {
		defer fileHandler.Close()
	}

	fileList := []*os.File{}

	renamedFilename, err := fileHandler.CreateFile(file.File{
		FileMeta: file.FileMeta{
			Name:      "rename",
			Extention: "json",
		},
		Content: "[]",
	})

	if err != nil {
		return closeFunc, "", err
	}

	fileList = append(fileList, fileHandler.GetFileWithLocked(renamedFilename))

	ctx, cancel := context.WithCancel(context.Background())
	workerManager := worker.NewWorkerManager[velog.VelogPostsItem, string](ctx, fmt.Sprintf("%s-%s", "velog-post-zip", username), 5)
	defer workerManager.Close()

	posts := getAllPosts(&velogApi)
	fileNameList := workerManager.Aggregate(cancel, posts,
		func(postItem velog.VelogPostsItem) string {
			post, err := velogApi.Post(postItem.UrlSlug)
			if err != nil {
				logger.Errorf("failed to get post %s, err: %s", post.Title, err.Error())
				return ""
			}

			sanitizedFile, isSanitized := sanitizeBasePathSpecialCase(post.Title)
			if isSanitized {
				logger.Debugf("sanitized file [%s] to [%s]", post.Title, sanitizedFile)
				err := fileHandler.AppendDataToJsonFile(renamedFilename, RenamedFileJSON{Origin: post.Title, Rename: sanitizedFile})
				if err != nil {
					logger.Error(err.Error())
				}
			}

			f, err := fileHandler.CreateFile(file.File{
				FileMeta: file.FileMeta{
					Name:      sanitizedFile,
					Extention: "md",
				},
				Content: post.Body,
			})

			if err != nil {
				logger.Error(err.Error())
			}
			return f
		})

	for _, filename := range fileNameList {
		if filename != "" {
			fileList = append(fileList, fileHandler.GetFileWithLocked(filename))
		}
	}

	zipFilename, err := fileHandler.CreateZipFile(file.ZipFile{
		FileMeta: file.FileMeta{
			Name:      fmt.Sprintf("%s-velog-posts", username),
			Extention: "zip",
		},
		Files: fileList,
	})

	if err != nil {
		return closeFunc, "", err
	}

	return closeFunc, zipFilename, nil
}

func getAllPosts(velogApi *velog.VelogAPI) velog.VelogPostsItemList {
	logger := log.GetLogger()
	velogPosts := velog.VelogPostsItemList{}
	cursor := ""
	for {
		posts, err := velogApi.Posts(cursor, 50)
		velogPosts = append(velogPosts, posts...)
		if err != nil {
			logger.Errorf("failed to get posts: %s", err)
			break
		}

		if len(posts) == 0 {
			break
		}

		lastPost := posts[len(posts)-1]
		post, err := velogApi.Post(lastPost.UrlSlug)
		if err != nil {
			logger.Errorf("failed to get post: %s", err)
			break
		}
		cursor = post.ID
	}
	return velogPosts
}

func FetchSelectedVelogPostsZip(username string, urlSlugList []string) (closeFunc, string, error) {
	// logger := log.GetLogger()
	velogApi := velog.NewVelogAPI(config.Manager.VelogConfig.ApiUrl, username)
	fh := file.NewFileHandler()

	closeFunc := func() {
		defer fh.Close()
	}

	fileList := []*os.File{}
	for _, urlSlug := range urlSlugList {
		post, err := velogApi.Post(urlSlug)

		// post에 데이터가 없는 경우 체크
		// 실패해도 계속 다운로드 해야하는 지
		if err != nil {
			return closeFunc, "", err
		}

		postFilePath, err := fh.CreateFile(file.File{
			FileMeta: file.FileMeta{
				Name:      post.Title,
				Extention: "md",
			},
			Content: post.Body,
		})

		if err != nil {
			return closeFunc, "", err
		}

		fileList = append(fileList, fh.GetFileWithLocked(postFilePath))
	}

	zipFilename, err := fh.CreateZipFile(file.ZipFile{
		FileMeta: file.FileMeta{
			Name:      fmt.Sprintf("%s-velog-posts", username),
			Extention: "zip",
		},
		Files: fileList,
	})

	if err != nil {
		return closeFunc, "", err
	}

	return closeFunc, zipFilename, nil
}
