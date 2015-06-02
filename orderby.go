package main

import (
	// "path/filepath"
	// "os"
	"sort"
)

type Files []Entry

func (self Files) Len() int           { return len(self) }
func (self Files) Less(i, j int) bool { return self[i].newpath < self[j].newpath }
func (self Files) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }

func (files Files) Sort() { sort.Sort(files) }

type ReverseSort struct{ Files }

func (self ReverseSort) Less(i, j int) bool { return self.Files[i].newpath > self.Files[j].newpath }

type MtimeSort struct{ Files }

func (self MtimeSort) Less(i, j int) bool {
	return self.Files[i].info.ModTime().UnixNano() < self.Files[j].info.ModTime().UnixNano()
}

type MtimeReverseSort struct{ Files }

func (self MtimeReverseSort) Less(i, j int) bool {
	return self.Files[i].info.ModTime().UnixNano() > self.Files[j].info.ModTime().UnixNano()
}

type SizeSort struct{ Files }

func (self SizeSort) Less(i, j int) bool {
	return self.Files[i].info.Size() > self.Files[i].info.Size()
}

type SizeReverseSort struct{ Files }

func (self SizeReverseSort) Less(i, j int) bool {
	return self.Files[i].info.Size() < self.Files[i].info.Size()
}
