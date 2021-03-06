package fsrename

import "os"

// The rename action
type Rename struct {
	dryrun bool
}

func (r *Rename) Act(entry *FileEntry) bool {
	if entry.newpath == "" || entry.newpath == entry.path {
		entry.message = "unchanged"
		return false
	}

	stat, err := os.Stat(entry.newpath)
	if os.IsExist(err) {
		entry.message = "file exists, ignore"
		return false
	}
	_ = stat

	if r.dryrun == false {
		os.Rename(entry.path, entry.newpath)
	}
	entry.message = "success"
	return true
}

func NewRenamer(dryrun bool) *Actor {
	return NewActor(&Rename{dryrun})
}
