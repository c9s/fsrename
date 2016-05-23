package fsrename

import "strings"
import "regexp"
import "fmt"

type Action interface {
	Act(entry *FileEntry) bool
}

type StrFormatReplaceAction struct {
	Search        string
	ReplaceFormat string
	N             int
	Seq           *Sequence
}

func NewStrFormatReplaceAction(search, replaceFormat string) *StrFormatReplaceAction {
	return &StrFormatReplaceAction{search, replaceFormat, 1, NewSequence(0, 1)}
}

func (s *StrFormatReplaceAction) Act(entry *FileEntry) bool {
	format := fmt.Sprintf(s.ReplaceFormat, s.Seq.Next())
	entry.newpath = strings.Replace(entry.path, s.Search, format, s.N)
	return true
}

type StrReplaceAction struct {
	Search  string
	Replace string
	N       int
}

func NewStrReplaceAction(search, replace string, n int) *StrReplaceAction {
	return &StrReplaceAction{search, replace, n}
}

func (s *StrReplaceAction) Act(entry *FileEntry) bool {
	entry.newpath = strings.Replace(entry.path, s.Search, s.Replace, s.N)
	return true
}

type RegExpAction struct {
	Matcher *regexp.Regexp
	Replace string
}

func NewRegExpAction(matcher *regexp.Regexp, replace string) *RegExpAction {
	return &RegExpAction{matcher, replace}
}

func NewRegExpActionWithPattern(pattern string, replace string) *RegExpAction {
	matcher := regexp.MustCompile(pattern)
	return &RegExpAction{matcher, replace}
}

func (s *RegExpAction) Act(entry *FileEntry) bool {
	entry.newpath = s.Matcher.ReplaceAllString(entry.path, s.Replace)
	return true
}
