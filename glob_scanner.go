package fsrename

import "path/filepath"
import "os"

type GlobScanner struct {
	*BaseWorker
}

func NewGlobScanner() *GlobScanner {
	return &GlobScanner{NewBaseWorker()}
}

func (self *GlobScanner) Run() {
	for {
		select {
		case <-self.stop:
			return
			break
		case entry := <-self.input:
			// end of data
			if entry == nil {
				self.emitEnd()
				return
			}
			matches, err := filepath.Glob(entry.path)
			if err != nil {
				return
			}
			for _, match := range matches {
				var err = filepath.Walk(match, func(path string, info os.FileInfo, err error) error {
					self.output <- NewFileEntryWithInfo(path, info)
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
