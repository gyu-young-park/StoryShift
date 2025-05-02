package service

import (
	"fmt"
	"os"

	"github.com/gyu-young-park/StoryShift/internal/config"
	"github.com/gyu-young-park/StoryShift/pkg/file"
	"github.com/gyu-young-park/StoryShift/pkg/log"
	"github.com/gyu-young-park/StoryShift/pkg/velog"
)

func GetPost(username, urlSlug string) (velog.VelogPost, error) {
	velogApi := velog.NewVelogAPI(config.Manager.VelogConfig.URL, username)
	velogPost, err := velogApi.Post(urlSlug)
	if err != nil {
		return velog.VelogPost{}, err
	}

	return velogPost, nil
}

func GetPosts(username, postId string, count int) ([]velog.VelogPostsItem, error) {
	velogApi := velog.NewVelogAPI(config.Manager.VelogConfig.URL, username)
	posts, err := velogApi.Posts(postId, count)
	if err != nil {
		return []velog.VelogPostsItem{}, err
	}

	return posts, nil
}

type closeFunc func()

func FetchVelogPostZip(username, postId string) (closeFunc, string, error) {
	velogApi := velog.NewVelogAPI(config.Manager.VelogConfig.URL, username)
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

func FetchAllVelogPostsZip(username string) (closeFunc, string, error) {
	logger := log.GetLogger()
	velogApi := velog.NewVelogAPI(config.Manager.VelogConfig.URL, username)

	fileHandler := file.NewFileHandler()
	closeFunc := func() {
		defer fileHandler.Close()
	}
	fileList := []*os.File{}
	cursor := ""

	for {
		posts, err := velogApi.Posts(cursor, 10)
		if err != nil {
			return closeFunc, "", err
		}

		if len(posts) == 0 {
			break
		}

		for _, postInfo := range posts {
			post, err := velogApi.Post(postInfo.UrlSlug)
			if err != nil {
				return closeFunc, "", err
			}

			f, err := fileHandler.CreateFile(file.File{
				FileMeta: file.FileMeta{
					Name:      sanitizeBasePathSpecialCase(post.Title),
					Extention: "md",
				},
				Content: post.Body,
			})

			if err != nil {
				logger.Error(err.Error())
			}

			fileList = append(fileList, fileHandler.GetFile(f))
			cursor = postInfo.ID
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
