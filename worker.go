package fsrename

// CondVar is used in Done, Stop signals
type CondVar chan bool

// Worker interface defines the implementation for
// streaming the file entries.
type Worker interface {
	// Set the input stream
	SetInput(input FileStream)

	// Set the output stream
	SetOutput(output FileStream)

	// Get the input stream
	Input() FileStream

	// Get the output stream
	Output() FileStream

	// Send stop signal to the listener go routine
	Stop()

	// Start listening the input channel
	Run()

	// Append a next worker to the next worker list
	// Returns the connected worker object
	Chain(nextWorker Worker) Worker
}

// Base Worker implements the common methods of a worker
type BaseWorker struct {
	input       FileStream
	output      FileStream
	stop        CondVar
	nextWorkers []Worker
}

// NewBaseWorker creates a base worker object with a default stop signal
// channel
func NewBaseWorker() *BaseWorker {
	return &BaseWorker{stop: make(CondVar, 1)}
}

// Stop a worker
func (w *BaseWorker) Stop() {
	// send stop signal to self
	w.stop <- true
	if len(w.nextWorkers) > 0 {
		for _, worker := range w.nextWorkers {
			worker.Stop()
		}
	}
}

// emitEnd sends nil to the next workers
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

// Chain() connects the next worker to children workers
// setup the parent output channel and the child input channel.
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

// Buffer all entries for further operations
// Returns FileEntries
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
	for i, _ := range entries {
		w.output <- &entries[i]
	}
}
