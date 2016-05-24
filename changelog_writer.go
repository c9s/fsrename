package fsrename

import (
	"encoding/csv"
	"log"
	"os"
)

type ChangeLogWriter struct {
	*BaseWorker
	logfile string
}

func NewChangeLogWriter(logfile string) *ChangeLogWriter {
	return &ChangeLogWriter{NewBaseWorker(), logfile}
}

func (w *ChangeLogWriter) Start() {
	go w.Run()
}

func (w *ChangeLogWriter) Run() {
	f, err := os.Create(w.logfile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	for {
		select {
		case <-w.stop:
			return
		case entry := <-w.input:
			// end of data
			if entry == nil {
				w.emitEnd()
				return
			}
			if err := writer.Write([]string{entry.path, entry.newpath, entry.message}); err != nil {
				log.Fatalln("Error writing files to csv:", err)
			}
		}
	}
	if err := writer.Error(); err != nil {
		log.Fatalln("Error writing files to csv", err)
	}
}
