package fsrename

import "regexp"

type RegExpFilter struct {
	*BaseWorker
	Matcher *regexp.Regexp
}

func NewRegExpFilter(matcher *regexp.Regexp) *RegExpFilter {
	return &RegExpFilter{NewBaseWorker(), matcher}
}

func NewRegExpFilterWithPattern(pattern string) *RegExpFilter {
	matcher := regexp.MustCompile(pattern)
	return &RegExpFilter{NewBaseWorker(), matcher}
}

func (w *RegExpFilter) Run() {
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
			if w.Matcher.MatchString(entry.info.Name()) {
				w.output <- entry
			}
		}
	}
}
