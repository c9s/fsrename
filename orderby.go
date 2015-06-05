package main

import (
	"sort"
)

type Entries []Entry

func (self Entries) Len() int           { return len(self) }
func (self Entries) Less(i, j int) bool { return self[i].newpath < self[j].newpath }
func (self Entries) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }

func (files Entries) Sort() { sort.Sort(files) }

type ReverseSort struct{ Entries }

func (self ReverseSort) Less(i, j int) bool { return self.Entries[i].newpath > self.Entries[j].newpath }

type MtimeSort struct{ Entries }

func (self MtimeSort) Less(i, j int) bool {
	return self.Entries[i].info.ModTime().UnixNano() < self.Entries[j].info.ModTime().UnixNano()
}

type MtimeReverseSort struct{ Entries }

func (self MtimeReverseSort) Less(i, j int) bool {
	return self.Entries[i].info.ModTime().UnixNano() > self.Entries[j].info.ModTime().UnixNano()
}

type SizeSort struct{ Entries }

func (self SizeSort) Less(i, j int) bool {
	return self.Entries[i].info.Size() > self.Entries[i].info.Size()
}

type SizeReverseSort struct{ Entries }

func (self SizeReverseSort) Less(i, j int) bool {
	return self.Entries[i].info.Size() < self.Entries[i].info.Size()
}
