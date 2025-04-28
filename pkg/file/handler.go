package file

import (
	"archive/zip"
	"fmt"
	"io"
	"os"

	"github.com/gyu-young-park/VelogStoryShift/pkg/log"
)

var Handler = fileHandler{}

type fileHandler struct {
}

type closeFileFunc func()

func (f fileHandler) CreateFile(file File) (closeFileFunc, *os.File, error) {
	logger := log.GetLogger()
	logger.Infof("download file: %s", file.GetFilename())

	dir := "./temp"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}

	tmpFile, err := os.Create(fmt.Sprintf("%s/%s", dir, file.GetFilename()))
	if err != nil {
		return func() {}, nil, err
	}
	logger.Infof("temp file: %s", tmpFile.Name())

	closeFunc := func() {
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()
	}

	if _, err := tmpFile.WriteString(file.Content); err != nil {

		return closeFunc, nil, err
	}

	if _, err := tmpFile.Seek(0, io.SeekStart); err != nil {
		os.Remove(tmpFile.Name())
		tmpFile.Close()
		return closeFunc, nil, err
	}

	return closeFunc, tmpFile, nil
}

func (f fileHandler) CreateZipFile(zipFileInfo ZipFile) (closeFileFunc, *os.File, error) {
	close, zipFile, err := f.CreateFile(File{
		FileMeta: FileMeta{
			Name:      "post",
			Extention: "zip",
		},
		Content: "",
	})

	closeFunc := func() {
		defer close()
	}

	if err != nil {
		return closeFunc, nil, err
	}

	zipWriter := zip.NewWriter(zipFile)
	zipWriter.Close()

	for _, file := range zipFileInfo.Files {
		zip, err := zipWriter.Create(file.Name())
		if err != nil {
			return closeFunc, nil, err
		}

		if _, err := io.Copy(zip, file); err != nil {
			return closeFunc, nil, err
		}
	}

	return closeFunc, zipFile, nil
}

func (f fileHandler) DonwloadPost(file File) (closeFileFunc, *os.File, error) {
	logger := log.GetLogger()
	close, titleFile, err := f.CreateFile(File{
		FileMeta: file.FileMeta,
		Content:  file.FileMeta.Name,
	})

	defer close()
	if err != nil {
		logger.Errorf("failed to create subject file:%s", titleFile.Name())
		return func() {}, nil, err
	}

	close, contentFile, err := f.CreateFile(file)
	defer close()

	if err != nil {
		logger.Errorf("failed to create content file:%s", contentFile.Name())
		return func() {}, nil, err
	}

	closeZipFile, zip, err := f.CreateZipFile(ZipFile{
		FileMeta: FileMeta{
			Name:      "post",
			Extention: "zip",
		},
		Files: []*os.File{
			titleFile, contentFile,
		},
	})

	if err != nil {
		logger.Errorf("failed to create zip file:%s", err)
		return func() {
			defer closeZipFile()
		}, nil, err
	}

	return func() {
		defer closeZipFile()
	}, zip, nil
}
