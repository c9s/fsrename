package fsrename

import "os"
import "path"

type FileEntry struct {
	path    string
	dir     string
	base    string
	info    os.FileInfo
	newpath string
	message string
}

func MustNewFileEntry(filepath string) *FileEntry {
	entry, err := NewFileEntry(filepath)
	if err != nil {
		panic(err)
	}
	return entry
}

func NewFileEntryWithInfo(filepath string, info os.FileInfo) *FileEntry {
	return &FileEntry{
		path: filepath,
		info: info,
		base: path.Base(filepath),
		dir:  path.Dir(filepath),
	}
}

func NewFileEntry(filepath string) (*FileEntry, error) {
	info, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}
	e := FileEntry{
		path: filepath,
		info: info,
		base: path.Base(filepath),
		dir:  path.Dir(filepath),
	}
	return &e, nil
}
