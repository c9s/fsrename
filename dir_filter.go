package fsrename

type DirFilter struct {
	*BaseWorker
}

func NewDirFilter() *DirFilter {
	return &DirFilter{NewBaseWorker()}
}

func (w *DirFilter) Start() {
	go w.Run()
}

func (w *DirFilter) Run() {
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
			if entry.info.Mode().IsDir() {
				w.output <- entry
			}
		}
	}
}
