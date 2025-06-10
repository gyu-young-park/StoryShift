package file

import (
	"path/filepath"
	"strings"
)

func SplitFilenameWithNameAndExt(filename string) (string, string) {
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	ext := strings.TrimPrefix(filepath.Ext(filename), ".")
	return name, ext
}
