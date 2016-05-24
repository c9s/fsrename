package fsrename

type Actor struct {
	*BaseWorker
	Action Action
}

func NewActor(a Action) *Actor {
	return &Actor{NewBaseWorker(), a}
}
func (a *Actor) Start() {
	go a.Run()
}

func (a *Actor) Run() {
	for {
		select {
		case <-a.stop:
			return
		case entry := <-a.input:
			// end of data
			if entry == nil {
				a.emitEnd()
				return
			}
			a.Action.Act(entry)
			a.output <- entry
		}
	}
}
