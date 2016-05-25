package fsrename

type Proxy struct {
	*BaseWorker
}

func NewProxy() *Proxy {
	return &Proxy{NewBaseWorker()}
}

func (s *Proxy) Start() {
	go s.Run()
}

func (s *Proxy) Run() {
	for {
		select {
		case <-s.stop:
			return
			break
		case entry := <-s.input:
			// end of data
			if entry == nil {
				s.emitEnd()
				return
			}
			if entry.info == nil {
				break
			}
			s.output <- entry
		}
	}
}
