package file

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/gyu-young-park/StoryShift/pkg/log"
)

const (
	TEMP_DIR = "/tmp/StoryShift"
)

func NewFileHandler() *FileHandler {
	if _, err := os.Stat(TEMP_DIR); os.IsNotExist(err) {
		err := os.MkdirAll(TEMP_DIR, 0755)
		if err != nil {
			panic(fmt.Sprintf("failed to create %s directory: %v", TEMP_DIR, err))
		}
	}

	fmt.Println("tempDir" + TEMP_DIR)

	dir, err := os.MkdirTemp(TEMP_DIR, "filehandler-*")
	if err != nil {
		panic(err)
	}

	fmt.Println("Random dir:", dir)

	return &FileHandler{
		mu:    sync.RWMutex{},
		dir:   dir,
		files: make(map[string]*os.File),
	}
}

type FileHandler struct {
	mu    sync.RWMutex
	dir   string
	files map[string]*os.File
}

func (f *FileHandler) GetFileWithLocked(filename string) *os.File {
	f.mu.Lock()
	defer f.mu.Unlock()

	logger := log.GetLogger()
	fh, ok := f.files[filename]
	if !ok {
		logger.Errorf("there is no file in this file handler: %s", filename)
		return nil
	}
	return fh
}

func (f *FileHandler) SetFileWithLocked(filename string, file *os.File) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.files[filename] = file
}

func (f *FileHandler) Close() {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, file := range f.files {
		os.Remove(file.Name())
		file.Close()
	}
	os.Remove(f.dir)
}

func (f *FileHandler) CreateFile(file File) (string, error) {
	logger := log.GetLogger()
	logger.Infof("download file: %s", file.GetFilename())

	tmpFile, err := os.Create(fmt.Sprintf("%s/%s", f.dir, file.GetFilename()))
	if err != nil {
		logger.Errorf("failed to create temp file: %s", file.GetFilename())
		return "", err
	}

	f.SetFileWithLocked(tmpFile.Name(), tmpFile)
	logger.Infof("temp file: %s", tmpFile.Name())

	if _, err := tmpFile.WriteString(file.Content); err != nil {
		logger.Errorf("failed to write content in temp file: %s", file.GetFilename())
		return "", err
	}

	if _, err := tmpFile.Seek(0, io.SeekStart); err != nil {
		logger.Errorf("failed to put file pointer on start point of temp file: %s", file.GetFilename())
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

	zipWriter := zip.NewWriter(f.GetFileWithLocked(zipFile))
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
			f.GetFileWithLocked(titleFile), f.GetFileWithLocked(contentFile),
		},
	})

	if err != nil {
		logger.Errorf("failed to create zip file:%s", err)
		return "", err
	}

	return zip, nil
}

func (f *FileHandler) AppendDataToJsonFile(filename string, data any) error {
	logger := log.GetLogger()
	fh := f.GetFileWithLocked(filename)
	if fh == nil {
		return fmt.Errorf("can't find file [%s] please create file", filename)
	}

	existingData := []any{}
	if err := json.NewDecoder(fh).Decode(&existingData); err != nil && err.Error() != "EOF" {
		logger.Errorf("failed to read existing data in [%s]", fh.Name())
		return err
	}

	existingData = append(existingData, data)

	fh.Truncate(0)
	fh.Seek(0, 0)

	if err := json.NewEncoder(fh).Encode(existingData); err != nil {
		logger.Errorf("failed to write existing data in [%s]", fh.Name())
		return err
	}

	fh.Seek(0, 0)

	return nil
}
