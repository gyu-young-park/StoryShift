package markdown

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/gyu-young-park/StoryShift/pkg/file"
	"github.com/gyu-young-park/StoryShift/pkg/log"
)

const (
	IMAGE_REGEX_MATCHER = `!\[[^\]]*\]\(([^)]+)\)`
	IMAGE_DIR_KEY       = "IMAGE_DIRECTORY_PATH_KEY"
)

type MarkdownImageManipulator struct {
	lock                             *sync.Mutex
	handler                          MarkdownImageHandlable
	downloadImageWithUrlReqModelList []DownloadImageWithUrlReqModel
	fileHandler                      *file.FileHandler
}

func NewMarkdownImageManipulator(handler MarkdownImageHandlable) *MarkdownImageManipulator {
	return &MarkdownImageManipulator{
		lock:                             &sync.Mutex{},
		handler:                          handler,
		downloadImageWithUrlReqModelList: []DownloadImageWithUrlReqModel{},
		fileHandler:                      file.NewFileHandler(),
	}
}

func (m *MarkdownImageManipulator) Replace(title string, contents string) string {
	m.lock.Lock()
	defer m.lock.Unlock()
	imageUrlList := m.handler.GetImageList(contents)
	if len(imageUrlList) == 0 {
		return contents
	}

	encodedTitle := base64.RawURLEncoding.EncodeToString([]byte(title))
	for i, image := range imageUrlList {
		m.downloadImageWithUrlReqModelList = append(m.downloadImageWithUrlReqModelList, DownloadImageWithUrlReqModel{
			Url:           image,
			ImageFileName: fmt.Sprintf("%s-%v", encodedTitle, i),
		})
	}

	return m.handler.ReplaceAllImageUrlOfContensWithPrefix(fmt.Sprintf("%s/%s", IMAGE_DIR_KEY, encodedTitle), contents)
}

func (m *MarkdownImageManipulator) DownloadAsZip(username string) (*os.File, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	logger := log.GetLogger()

	imageFileList := []*os.File{}
	downloadImageWithUrlRespModel, err := m.handler.DownloadImageWithUrl(m.fileHandler, m.downloadImageWithUrlReqModelList)
	if err != nil {
		logger.Errorf("failed to download %s", err.Error())
		return nil, err
	}

	if len(downloadImageWithUrlRespModel.FailedToDownloadImageUrlList) > 0 {
		bFailedImageList, _ := json.Marshal(downloadImageWithUrlRespModel.FailedToDownloadImageUrlList)
		downloadFailListFilePath, err := m.fileHandler.CreateFile(file.File{
			FileMeta: file.FileMeta{
				Name:      "download_fail_list",
				Extention: "json",
			},
			Content: string(bFailedImageList),
		})

		if err != nil {
			logger.Errorf("failed to create download fail list, err: %s", err.Error())
		} else {
			imageFileList = append(imageFileList, m.fileHandler.GetFileWithLocked(downloadFailListFilePath))
		}
	}

	for _, image := range downloadImageWithUrlRespModel.ImageFilePathList {
		if image != "" {
			imageFileList = append(imageFileList, m.fileHandler.GetFileWithLocked(image))
		}
	}

	imageZipname, err := m.fileHandler.CreateZipFile(file.ZipFile{
		FileMeta: file.FileMeta{
			Name:      fmt.Sprintf("%s-%s", username, "image"),
			Extention: "zip",
		},
		Files: imageFileList,
	})

	if err != nil {
		return nil, err
	}

	imageZipfile := m.fileHandler.GetFileWithLocked(imageZipname)
	imageZipfile.Seek(0, 0)

	return imageZipfile, nil
}

func (m *MarkdownImageManipulator) Done() {
	m.downloadImageWithUrlReqModelList = []DownloadImageWithUrlReqModel{}
	m.fileHandler.Close()
	m.fileHandler = nil
}
