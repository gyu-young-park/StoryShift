package markdown

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

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
	index := 0
	return m.matcher.ReplaceAllStringFunc(contents, func(match string) string {
		imageName := fmt.Sprintf("%v-%v", imageNamePrefix, index)
		index++
		submatch := m.matcher.FindStringSubmatch(match)
		urlIndex := len(submatch) - 1
		alt := match[:len(match)-len(submatch[urlIndex])-1]
		return alt + imageName + ")"
	})
}

func (m *MKImageHander) DownloadImageWithUrl(fh *file.FileHandler, imageUrls []string) []string {
	logger := log.GetLogger()

	if imageUrls == nil {
		logger.Error("there is no imageUrls map")
		return []string{}
	}

	imageFileList := []string{}
	for _, imageUrl := range imageUrls {
		resp, err := httpclient.Get(httpclient.GetRequestParam{
			URL: imageUrl,
		})

		if err != nil {
			logger.Errorf("failed to download image: %s", imageUrl)
			continue
		}

		imageNameAndExt := strings.Split(filepath.Base(imageUrl), ".")
		if len(imageNameAndExt) < 2 {
			logger.Errorf("invalid image url: %s", imageUrl)
			continue
		}

		imageName := imageNameAndExt[0]
		ext := imageNameAndExt[1]

		imageFilePath, err := fh.CreateFile(file.File{
			FileMeta: file.FileMeta{
				Name:      imageName,
				Extention: ext,
			},
			Content: string(resp.Body),
		})

		if err != nil {
			logger.Errorf("failed to create image file: %s", imageUrl)
			continue
		}
		imageFileList = append(imageFileList, imageFilePath)
	}

	return imageFileList
}
