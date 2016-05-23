package fsrename

import "sync/atomic"

type Sequence struct {
	Start int32
	Last  int32
	Step  int32
}

func NewSequence(start int32, step int32) *Sequence {
	return &Sequence{start, start, step}
}

func (s *Sequence) Reset() {
	atomic.StoreInt32(&s.Last, 0)
}

func (s *Sequence) Next() int32 {
	atomic.AddInt32(&s.Last, 1)
	return s.Last
}
