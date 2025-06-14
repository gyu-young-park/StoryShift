package markdown

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/gyu-young-park/StoryShift/internal/httpclient"
	"github.com/gyu-young-park/StoryShift/pkg/file"
	"github.com/gyu-young-park/StoryShift/pkg/log"
)

type MarkdownImageHandlable interface {
	GetImageList(contents string) []string
	ReplaceAllImageUrlOfContensWithPrefix(imageNamePrefix string, contents string) string
	DownloadImageWithUrl(fh *file.FileHandler, imageUrls map[string]string) ([]string, []string, error)
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

func (m *MarkdownImageHandler) DownloadImageWithUrl(fh *file.FileHandler, imageUrls map[string]string) ([]string, []string, error) {
	logger := log.GetLogger()
	if imageUrls == nil {
		return []string{}, []string{}, fmt.Errorf("there is no imageUrls map")
	}

	failedImageList := []string{}
	files := []file.File{}
	for replaceName, imageUrl := range imageUrls {
		resp, err := httpclient.Get(httpclient.GetRequestParam{
			URL: imageUrl,
		})
		if err != nil {
			failedImageList = append(failedImageList, imageUrl)
		}

		_, ext := file.SplitFilenameWithNameAndExt(imageUrl)
		files = append(files, file.File{
			FileMeta: file.FileMeta{
				Name:      replaceName,
				Extention: ext,
			},
			Content: string(resp.Body),
		})
	}

	imageFileList := []string{}
	for _, file := range files {
		imageFilePath, err := fh.CreateFile(file)
		if err != nil {
			return []string{}, failedImageList, fmt.Errorf("failed to create image file: %s", file.GetFilename())
		}
		imageFileList = append(imageFileList, imageFilePath)
		logger.Infof("image file created: %s", imageFilePath)
	}

	return imageFileList, failedImageList, nil
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
