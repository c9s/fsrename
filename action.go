package fsrename

import "strings"
import "regexp"
import "fmt"
import "path"

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
	if strings.Contains(entry.base, s.Search) {
		newbase := strings.Replace(entry.base, s.Search, s.Replace, s.N)
		entry.newpath = path.Join(entry.dir, newbase)
		return true
	}
	return false
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
	if strings.Contains(entry.base, s.Search) {
		format := fmt.Sprintf(s.ReplaceFormat, s.Seq.Next())
		newbase := strings.Replace(entry.base, s.Search, format, s.N)
		entry.newpath = path.Join(entry.dir, newbase)
		return true
	}
	return false
}

type RegExpReplaceAction struct {
	Matcher *regexp.Regexp
	Replace string
}

func NewRegExpReplaceAction(matcher *regexp.Regexp, replace string) *RegExpReplaceAction {
	return &RegExpReplaceAction{matcher, replace}
}

func NewRegExpReplaceActionWithPattern(pattern string, replace string) *RegExpReplaceAction {
	matcher := regexp.MustCompile(pattern)
	return &RegExpReplaceAction{matcher, replace}
}

func (s *RegExpReplaceAction) Act(entry *FileEntry) bool {
	if s.Matcher.MatchString(entry.base) {
		newbase := s.Matcher.ReplaceAllString(entry.base, s.Replace)
		entry.newpath = path.Join(entry.dir, newbase)
		return true
	}
	return false
}

type RegExpFormatReplaceAction struct {
	Matcher       *regexp.Regexp
	ReplaceFormat string
	Seq           *Sequence
}

func NewRegExpFormatReplaceAction(matcher *regexp.Regexp, replaceFormat string) *RegExpFormatReplaceAction {
	return &RegExpFormatReplaceAction{matcher, replaceFormat, NewSequence(0, 1)}
}

func (s *RegExpFormatReplaceAction) Act(entry *FileEntry) bool {
	if s.Matcher.MatchString(entry.base) {
		format := fmt.Sprintf(s.ReplaceFormat, s.Seq.Next())
		newbase := s.Matcher.ReplaceAllString(entry.base, format)
		entry.newpath = path.Join(entry.dir, newbase)
		return true
	}
	return false
}
