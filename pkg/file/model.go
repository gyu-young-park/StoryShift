package file

import (
	"fmt"
	"os"
)

type FileMeta struct {
	Name      string
	Extention string
}

func (f FileMeta) GetFilename() string {
	return fmt.Sprintf("%s.%s", f.Name, f.Extention)
}

type File struct {
	FileMeta
	Content string
}

type ZipFile struct {
	FileMeta
	Files []*os.File
}
