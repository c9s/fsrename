package fsrename

type CondVar chan bool

type FileStream chan *FileEntry

func NewFileStream() FileStream {
	return make(FileStream, 10)
}

type WorkerHandle func(input FileStream)

type Worker interface {
	SetInput(input FileStream)
	SetOutput(output FileStream)
	Input() FileStream
	Output() FileStream
	Stop()
	Run()
	Chain(worker Worker) Worker
}

type BaseWorker struct {
	input       FileStream
	output      FileStream
	stop        CondVar
	nextWorkers []Worker
}

func NewBaseWorker() *BaseWorker {
	return &BaseWorker{stop: make(CondVar, 1)}
}

/**
Stop a worker
*/
func (w *BaseWorker) Stop() {
	// send stop signal to self
	w.stop <- true
	if len(w.nextWorkers) > 0 {
		for _, worker := range w.nextWorkers {
			worker.Stop()
		}
	}
}

func (w *BaseWorker) emitEnd() {
	if len(w.nextWorkers) > 0 {
		for _, worker := range w.nextWorkers {
			if nin := worker.Input(); nin != nil {
				nin <- nil
			}
		}
	}
	w.output <- nil
}

func (w *BaseWorker) Input() FileStream {
	return w.input
}

func (w *BaseWorker) Output() FileStream {
	return w.output
}

func (w *BaseWorker) SetInput(input FileStream) {
	w.input = input
}

func (w *BaseWorker) SetOutput(output FileStream) {
	w.output = output
}

func (w *BaseWorker) Chain(nextWorker Worker) Worker {
	if w.output == nil {
		w.output = NewFileStream()
	}
	// override the input in the worker
	nextWorker.SetInput(w.output)

	if o := nextWorker.Output(); o == nil {
		nextWorker.SetOutput(NewFileStream())
	}
	w.nextWorkers = append(w.nextWorkers, nextWorker)
	return nextWorker
}

func (w *BaseWorker) BufferEntries() []FileEntry {
	var entries []FileEntry
	for {
		select {
		case <-w.stop:
			return entries
		case entry := <-w.input:
			// end of data
			if entry == nil {
				return entries
			}
			entries = append(entries, *entry)
		}
	}
	return entries
}

func (w *BaseWorker) FlushEntries(entries []FileEntry) {
	for _, entry := range entries {
		w.output <- &entry
	}
}
