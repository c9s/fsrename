package fsrename

type FileFilter struct {
	*BaseWorker
}

func NewFileFilter() *FileFilter {
	return &FileFilter{NewBaseWorker()}
}

func (f *FileFilter) Start() {
	go f.Run()
}

func (f *FileFilter) Run() {
	for {
		select {
		case <-f.stop:
			return
			break
		case entry := <-f.input:
			// end of data
			if entry == nil {
				f.emitEnd()
				return
			}
			if entry.info.Mode().IsRegular() {
				f.output <- entry
			}
		}
	}
}
