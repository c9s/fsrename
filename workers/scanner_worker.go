package workers

// import "regexp"
import "os"
import "path/filepath"

type FileEntry struct {
	path string
	info os.FileInfo
}

type ScannerSettings struct {
	FileOnly bool
	DirOnly  bool
	// MatchRegExp *regexp.Regexp
}

type Worker func(cv chan bool, input chan *FileEntry, output chan *FileEntry)

type Producer func(cv chan bool, output chan *FileEntry)

func CreateScannerWorker(scanPath string, settings ScannerSettings) Worker {
	return func(cv chan bool, input chan *FileEntry, output chan *FileEntry) {
		matches, err := filepath.Glob(scanPath)
		if err != nil {
			cv <- false
			return
		}
		for _, match := range matches {
			var err = filepath.Walk(match, func(path string, info os.FileInfo, err error) error {
				if settings.DirOnly {
					if info.IsDir() {
						output <- &FileEntry{path: path, info: info}
					}
				} else if settings.FileOnly {
					if !info.IsDir() {
						output <- &FileEntry{path: path, info: info}
					}
				} else {
					output <- &FileEntry{path: path, info: info}
				}
				return err
			})
			if err != nil {
				panic(err)
			}
		}
		output <- nil
		cv <- true
	}
}
