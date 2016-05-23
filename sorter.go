package fsrename

import (
	"sort"
)

type FileEntries []FileEntry

func (self FileEntries) Len() int           { return len(self) }
func (self FileEntries) Less(i, j int) bool { return self[i].newpath < self[j].newpath }
func (self FileEntries) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }

func (files FileEntries) Sort() { sort.Sort(files) }

type ReverseSort struct{ FileEntries }

func (self ReverseSort) Less(i, j int) bool {
	return self.FileEntries[i].newpath > self.FileEntries[j].newpath
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

func (s *ReverseSorter) Run() {
	var entries = s.BufferEntries()
	sorter := ReverseSort{entries}
	sorter.Sort()
	s.FlushEntries(entries)
	s.emitEnd()
}

type MtimeSorter struct {
	*BaseWorker
}

func (s *MtimeSorter) Run() {
	var entries = s.BufferEntries()
	sorter := MtimeSort{entries}
	sorter.Sort()
	s.FlushEntries(entries)
	s.emitEnd()
}

type MtimeReverseSorter struct {
	*BaseWorker
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

func (s *SizeReverseSorter) Run() {
	var entries = s.BufferEntries()
	sorter := SizeReverseSort{entries}
	sorter.Sort()
	s.FlushEntries(entries)
	s.emitEnd()
}
