package markdown

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/gyu-young-park/StoryShift/internal/httpclient"
	"github.com/gyu-young-park/StoryShift/pkg/file"
	"github.com/gyu-young-park/StoryShift/pkg/log"
	"github.com/gyu-young-park/StoryShift/pkg/worker"
)

type MarkdownImageHandlable interface {
	GetImageList(contents string) []string
	ReplaceAllImageUrlOfContensWithPrefix(imageNamePrefix string, contents string) string
	DownloadImageWithUrl(fh *file.FileHandler, imageUrls []DownloadImageWithUrlReqModel) (DownloadImageWithUrlRespModel, error)
}

type MarkdownImageHandler struct {
	matcher *regexp.Regexp
}

func NewMarkdownImageHandler() *MarkdownImageHandler {
	return &MarkdownImageHandler{
		matcher: regexp.MustCompile(IMAGE_REGEX_MATCHER),
	}
}

func (m *MarkdownImageHandler) GetImageList(contents string) []string {
	matches := m.matcher.FindAllStringSubmatch(contents, -1)

	images := []string{}
	for _, match := range matches {
		images = append(images, match[1])
	}

	return images
}

func (m *MarkdownImageHandler) ReplaceAllImageUrlOfContensWithPrefix(imageNamePrefix string, contents string) string {
	logger := log.GetLogger()
	index := 0
	return m.matcher.ReplaceAllStringFunc(contents, func(match string) string {
		imageName := fmt.Sprintf("%v-%v", imageNamePrefix, index)
		index++

		submatch := m.matcher.FindStringSubmatch(match)
		urlIndex := len(submatch) - 1
		logger.Debugf("change image url from [%s] to [%s]", submatch[urlIndex-1], imageName)
		alt := match[:len(match)-len(submatch[urlIndex])-1]
		return alt + imageName + filepath.Ext(match)
	})
}

func (m *MarkdownImageHandler) DownloadImageWithUrl(fh *file.FileHandler, requests []DownloadImageWithUrlReqModel) (DownloadImageWithUrlRespModel, error) {
	logger := log.GetLogger()
	if len(requests) == 0 {
		return DownloadImageWithUrlRespModel{}, fmt.Errorf("there is no req")
	}

	mu := sync.Mutex{}
	failedImageList := []string{}

	ctx, cancel := context.WithCancel(context.Background())
	workerManager := worker.NewWorkerManager[DownloadImageWithUrlReqModel, file.File](ctx, "markdown-image-downloader", 50)
	defer workerManager.Close()

	imageFileList := workerManager.Aggregate(cancel, requests, func(req DownloadImageWithUrlReqModel) file.File {
		resp, err := httpclient.Get(httpclient.GetRequestParam{
			URL: req.Url,
		})
		if err != nil {
			mu.Lock()
			defer mu.Unlock()
			failedImageList = append(failedImageList, req.Url)
			return file.File{
				FileMeta: file.FileMeta{
					Name:      "failed",
					Extention: ".fail",
				},
			}
		}

		_, ext := file.SplitFilenameWithNameAndExt(req.Url)
		return file.File{
			FileMeta: file.FileMeta{
				Name:      req.ImageFileName,
				Extention: ext,
			},
			Content: string(resp.Body),
		}
	})

	imageFilePathList := []string{}
	for _, file := range imageFileList {
		if file.Name == "failed" {
			continue
		}

		imageFilePath, err := fh.CreateFile(file)
		if err != nil {
			return DownloadImageWithUrlRespModel{}, fmt.Errorf("failed to create image file: %s", file.GetFilename())
		}
		imageFilePathList = append(imageFilePathList, imageFilePath)
		logger.Infof("image file created: %s", imageFilePath)
	}

	return DownloadImageWithUrlRespModel{
		ImageFilePathList:            imageFilePathList,
		FailedToDownloadImageUrlList: failedImageList}, nil
}

type DefaultMarkdownImageHandler struct {
}

func NewDefaultMarkdownImageHandler() *DefaultMarkdownImageHandler {
	return &DefaultMarkdownImageHandler{}
}

func (m *DefaultMarkdownImageHandler) GetImageList(contents string) []string {
	return []string{}
}

func (m *DefaultMarkdownImageHandler) ReplaceAllImageUrlOfContensWithPrefix(imageNamePrefix string, contents string) string {
	return ""
}

func (m *DefaultMarkdownImageHandler) DownloadImageWithUrl(fh *file.FileHandler, imageUrls map[string]string) ([]string, error) {
	return []string{}, nil
}
