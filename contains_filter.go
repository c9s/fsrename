package fsrename

import "strings"

type StrContainsFilter struct {
	*BaseWorker
	Search string
}

func NewStrContainsFilter(search string) *StrContainsFilter {
	return &StrContainsFilter{NewBaseWorker(), search}
}

func (w *StrContainsFilter) Start() {
	go w.Run()
}

func (w *StrContainsFilter) Run() {
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
			if strings.Contains(entry.base, w.Search) {
				w.output <- entry
			}
		}
	}
}
