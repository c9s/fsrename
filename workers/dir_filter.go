package workers

type DirFilter struct {
	*BaseWorker
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
			if entry.info.IsDir() {
				w.output <- entry
			}
		}
	}
}
