package workers

type Actor struct {
	*BaseWorker
	Action Action
}

func NewActor(a Action) *Actor {
	return &Actor{NewBaseWorker(), a}
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
		}
	}
}
