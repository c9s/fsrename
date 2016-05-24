package fsrename

import "path/filepath"
import "os"

type GlobScanner struct {
	*BaseWorker
}

func NewGlobScanner() *GlobScanner {
	return &GlobScanner{NewBaseWorker()}
}

func (s *GlobScanner) Start() {
	go s.Run()
}

func (s *GlobScanner) Run() {
	for {
		select {
		case <-s.stop:
			return
			break
		case entry := <-s.input:
			// end of data
			if entry == nil {
				s.emitEnd()
				return
			}
			matches, err := filepath.Glob(entry.path)
			if err != nil {
				return
			}
			for _, match := range matches {
				var err = filepath.Walk(match, func(path string, info os.FileInfo, err error) error {
					base := filepath.Base(path)
					if base == ".svn" || base == ".git" || base == ".hg" {
						return filepath.SkipDir
					}
					if err != nil {
						panic(err)
					}
					s.output <- NewFileEntryWithInfo(path, info)
					return err
				})
				if err != nil {
					panic(err)
				}
			}
			break
		}
	}
}
