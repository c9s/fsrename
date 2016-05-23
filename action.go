package workers

import "strings"
import "regexp"

type Action interface {
	Act(entry *FileEntry) bool
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
