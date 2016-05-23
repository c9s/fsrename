package workers

import "os"

type FileEntry struct {
	path    string
	info    os.FileInfo
	newpath string

	// an object pointer with Action interface
	action Action
}

func MustNewFileEntry(filepath string) *FileEntry {
	info, err := os.Stat(filepath)
	if err != nil {
		panic(err)
	}
	return &FileEntry{filepath, info, "", nil}
}

func NewFileEntry(filepath string) (*FileEntry, error) {
	info, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}
	e := FileEntry{filepath, info, "", nil}
	return &e, nil
}
