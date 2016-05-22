package workers

import "os"

type FileEntry struct {
	path    string
	info    os.FileInfo
	newpath string
}

func NewFileEntry(filepath string) (*FileEntry, error) {
	info, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}
	e := FileEntry{filepath, info, ""}
	return &e, nil
}
