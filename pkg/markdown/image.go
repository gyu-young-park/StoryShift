package markdown

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/gyu-young-park/StoryShift/internal/httpclient"
	"github.com/gyu-young-park/StoryShift/pkg/file"
	"github.com/gyu-young-park/StoryShift/pkg/log"
)

const (
	IMAGE_REGEX_MATCHER = `!\[[^\]]*\]\(([^)]+)\)`
)

type MKImageHander struct {
	matcher *regexp.Regexp
}

func NewMKImageHander() *MKImageHander {
	return &MKImageHander{
		matcher: regexp.MustCompile(IMAGE_REGEX_MATCHER),
	}
}

func (m *MKImageHander) GetImageList(contents string) []string {
	matches := m.matcher.FindAllStringSubmatch(contents, -1)

	images := []string{}
	for _, match := range matches {
		images = append(images, match[1])
	}

	return images
}

func (m *MKImageHander) ReplaceAllImageUrlOfContensWithPrefix(imageNamePrefix string, contents string) string {
	logger := log.GetLogger()
	index := 0
	return m.matcher.ReplaceAllStringFunc(contents, func(match string) string {
		imageName := fmt.Sprintf("%v-%v", imageNamePrefix, index)
		index++

		submatch := m.matcher.FindStringSubmatch(match)
		urlIndex := len(submatch) - 1
		logger.Debugf("change image url from [%s] to [%s]", submatch[urlIndex-1], imageName)
		alt := match[:len(match)-len(submatch[urlIndex])-1]
		return alt + imageName + filepath.Ext(match) + ")"
	})
}

func (m *MKImageHander) DownloadImageWithUrl(fh *file.FileHandler, imageUrls map[string]string) ([]string, error) {
	logger := log.GetLogger()

	if imageUrls == nil {
		logger.Error("there is no imageUrls map")
		return []string{}, fmt.Errorf("there is no imageUrls map")
	}

	files := []file.File{}
	for replaceName, imageUrl := range imageUrls {
		resp, err := httpclient.Get(httpclient.GetRequestParam{
			URL: imageUrl,
		})

		if err != nil {
			logger.Errorf("failed to download image: %s", imageUrl)
			return []string{}, fmt.Errorf("failed to download image: %s, error: %v", imageUrl, err)
		}

		filename, ext := file.SplitFilenameWithNameAndExt(replaceName)
		files = append(files, file.File{
			FileMeta: file.FileMeta{
				Name:      filename,
				Extention: ext,
			},
			Content: string(resp.Body),
		})

		if err != nil {
			logger.Errorf("failed to create image file: %s", imageUrl)
			continue
		}
	}

	imageFileList := []string{}
	for _, file := range files {
		imageFilePath, err := fh.CreateFile(file)
		if err != nil {
			logger.Errorf("failed to create image file: %s", file.GetFilename())
			return []string{}, fmt.Errorf("failed to create image file: %s", file.GetFilename())
		}
		imageFileList = append(imageFileList, imageFilePath)
		logger.Infof("image file created: %s", imageFilePath)
	}

	return imageFileList, nil
}
