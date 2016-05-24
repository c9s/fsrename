package fsrename

import (
	"sort"
)

func (files FileEntries) Len() int           { return len(files) }
func (files FileEntries) Less(i, j int) bool { return files[i].path < files[j].path }
func (files FileEntries) Swap(i, j int) {
	files[i], files[j] = files[j], files[i]
}
func (files FileEntries) Sort() { sort.Sort(files) }

type ReverseSort struct{ FileEntries }

func (self ReverseSort) Less(i, j int) bool {
	return self.FileEntries[i].path > self.FileEntries[j].path
}

type MtimeSort struct{ FileEntries }

func (self MtimeSort) Less(i, j int) bool {
	return self.FileEntries[i].info.ModTime().UnixNano() < self.FileEntries[j].info.ModTime().UnixNano()
}

type MtimeReverseSort struct{ FileEntries }

func (self MtimeReverseSort) Less(i, j int) bool {
	return self.FileEntries[i].info.ModTime().UnixNano() > self.FileEntries[j].info.ModTime().UnixNano()
}

type SizeSort struct{ FileEntries }

func (self SizeSort) Less(i, j int) bool {
	return self.FileEntries[i].info.Size() > self.FileEntries[i].info.Size()
}

type SizeReverseSort struct{ FileEntries }

func (self SizeReverseSort) Less(i, j int) bool {
	return self.FileEntries[i].info.Size() < self.FileEntries[i].info.Size()
}

type ReverseSorter struct {
	*BaseWorker
}

func NewReverseSorter() *ReverseSorter {
	return &ReverseSorter{NewBaseWorker()}
}

func (s *ReverseSorter) Run() {
	var entries = s.BufferEntries()
	sorter := ReverseSort{entries}
	sorter.Sort()
	s.FlushEntries(entries)
	s.emitEnd()
}

func (s *ReverseSorter) Start() {
	go s.Run()
}

type MtimeSorter struct {
	*BaseWorker
}

func NewMtimeSorter() *MtimeSorter {
	return &MtimeSorter{NewBaseWorker()}
}

func (s *MtimeSorter) Run() {
	var entries = s.BufferEntries()
	sorter := MtimeSort{entries}
	sorter.Sort()
	s.FlushEntries(entries)
	s.emitEnd()
}

func (s *MtimeSorter) Start() {
	go s.Run()
}

type MtimeReverseSorter struct {
	*BaseWorker
}

func NewMtimeReverseSorter() *MtimeReverseSorter {
	return &MtimeReverseSorter{NewBaseWorker()}
}

func (s *MtimeReverseSorter) Start() {
	go s.Run()
}

func (s *MtimeReverseSorter) Run() {
	var entries = s.BufferEntries()
	sorter := MtimeReverseSort{entries}
	sorter.Sort()
	s.FlushEntries(entries)
	s.emitEnd()
}

type SizeSorter struct {
	*BaseWorker
}

func NewSizeSorter() *SizeSorter {
	return &SizeSorter{NewBaseWorker()}
}

func (s *SizeSorter) Start() {
	go s.Run()
}

func (s *SizeSorter) Run() {
	var entries = s.BufferEntries()
	sorter := SizeSort{entries}
	sorter.Sort()
	s.FlushEntries(entries)
	s.emitEnd()
}

type SizeReverseSorter struct {
	*BaseWorker
}

func NewSizeReverseSorter() *SizeReverseSorter {
	return &SizeReverseSorter{NewBaseWorker()}
}

func (s *SizeReverseSorter) Start() {
	go s.Run()
}

func (s *SizeReverseSorter) Run() {
	var entries = s.BufferEntries()
	sorter := SizeReverseSort{entries}
	sorter.Sort()
	s.FlushEntries(entries)
	s.emitEnd()
}
