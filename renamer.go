package fsrename

import "os"

// The rename action
type Rename struct{}

func (r *Rename) Act(entry *FileEntry) bool {
	if entry.newpath == "" {
		entry.message = "empty file"
		return false
	}

	stat, err := os.Stat(entry.newpath)
	if os.IsExist(err) {
		entry.message = "file exists, ignore"
		return false
	}
	_ = stat

	os.Rename(entry.path, entry.newpath)
	entry.message = "success"
	return true
}

func NewRenamer() *Actor {
	return NewActor(&Rename{})
}
