package file

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gyu-young-park/VelogStoryShift/pkg/log"
)

const (
	TEMP_DIR = "./temp"
)

func NewFileHandler() *FileHandler {
	return &FileHandler{
		files: make(map[string]*os.File),
	}
}

type FileHandler struct {
	files map[string]*os.File
}

func (f *FileHandler) GetFile(filename string) *os.File {
	logger := log.GetLogger()
	fh, ok := f.files[filename]
	if !ok {
		logger.Errorf("there is no file in this file handler: %s", filename)
		return nil
	}
	return fh
}

func (f *FileHandler) Close() {
	for _, file := range f.files {
		os.Remove(file.Name())
		file.Close()
	}
}

func (f *FileHandler) CreateFile(file File) (string, error) {
	logger := log.GetLogger()
	logger.Infof("download file: %s", file.GetFilename())

	if _, err := os.Stat(TEMP_DIR); os.IsNotExist(err) {
		os.Mkdir(TEMP_DIR, 0755)
	}

	tmpFile, err := os.Create(fmt.Sprintf("%s/%s", TEMP_DIR, file.GetFilename()))
	if err != nil {
		return "", err
	}

	f.files[tmpFile.Name()] = tmpFile

	logger.Infof("temp file: %s", tmpFile.Name())

	if _, err := tmpFile.WriteString(file.Content); err != nil {
		return "", err
	}

	if _, err := tmpFile.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func (f *FileHandler) CreateZipFile(zipFileInfo ZipFile) (string, error) {
	zipFile, err := f.CreateFile(File{
		FileMeta: FileMeta{
			Name:      zipFileInfo.Name,
			Extention: "zip",
		},
		Content: "",
	})

	if err != nil {
		return "", err
	}

	zipWriter := zip.NewWriter(f.files[zipFile])
	defer zipWriter.Close()

	for _, file := range zipFileInfo.Files {
		w, err := zipWriter.Create(filepath.Base(file.Name()))
		if err != nil {
			return "", err
		}

		if _, err := io.Copy(w, file); err != nil {
			return "", err
		}
	}

	return zipFile, nil
}

func (f *FileHandler) MakeZipFileWithTitle(file File) (string, error) {
	logger := log.GetLogger()
	titleFile, err := f.CreateFile(File{
		FileMeta: FileMeta{
			Name:      "title",
			Extention: "txt",
		},
		Content: file.FileMeta.Name,
	})

	if err != nil {
		logger.Errorf("failed to create subject file:%s", titleFile)
		return "", err
	}

	contentFile, err := f.CreateFile(file)

	if err != nil {
		logger.Errorf("failed to create content file:%s", contentFile)
		return "", err
	}

	zip, err := f.CreateZipFile(ZipFile{
		FileMeta: FileMeta{
			Name:      file.Name,
			Extention: "zip",
		},
		Files: []*os.File{
			f.files[titleFile], f.files[contentFile],
		},
	})

	if err != nil {
		logger.Errorf("failed to create zip file:%s", err)
		return "", err
	}

	return zip, nil
}
