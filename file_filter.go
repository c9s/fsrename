package workers

type FileFilter struct {
	*BaseWorker
}

func NewFileFilter() *FileFilter {
	return &FileFilter{NewBaseWorker()}
}

func (w *FileFilter) Run() {
	for {
		select {
		case <-w.stop:
			return
			break
		case entry := <-w.input:
			// end of data
			if entry == nil {
				w.emitEnd()
				return
			}
			if !entry.info.IsDir() {
				w.output <- entry
			}
		}
	}
}
