package fsrename

import "strings"
import "regexp"
import "fmt"
import "path"
import "strconv"

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

type UnderscoreAction struct {
	spliter *regexp.Regexp
}

func NewUnderscoreAction() *UnderscoreAction {
	spliter := regexp.MustCompile("[A-Z]")
	return &UnderscoreAction{spliter}
}

func (a *UnderscoreAction) Act(entry *FileEntry) bool {
	newbase := a.spliter.ReplaceAllStringFunc(entry.base, func(w string) string {
		return "_" + strings.ToLower(w)
	})
	newbase = strings.TrimLeft(newbase, "_")
	entry.newpath = path.Join(entry.dir, newbase)
	return true
}

type CamelCaseAction struct {
	spliter *regexp.Regexp
}

func NewCamelCaseAction(splitstr string) *CamelCaseAction {
	spliter := regexp.MustCompile(splitstr)
	return &CamelCaseAction{spliter}
}

func (a *CamelCaseAction) Act(entry *FileEntry) bool {
	substrings := a.spliter.Split(entry.base, -1)
	for i, str := range substrings {
		substrings[i] = strings.ToUpper(str[:1]) + str[1:]
	}
	entry.newpath = path.Join(entry.dir, strings.Join(substrings, ""))
	return true
}

// PrefixAction adds prefix to the matched filenames
type PrefixAction struct {
	Prefix string
}

func NewPrefixAction(prefix string) *PrefixAction {
	return &PrefixAction{prefix}
}

func (a *PrefixAction) Act(entry *FileEntry) bool {
	if strings.HasPrefix(entry.base, a.Prefix) {
		return false
	}
	entry.newpath = path.Join(entry.dir, a.Prefix+entry.base)
	return true
}

// SuffixAction adds prefix to the matched filenames
type SuffixAction struct {
	Suffix string
}

func NewSuffixAction(suffix string) *SuffixAction {
	return &SuffixAction{suffix}
}

func (a *SuffixAction) Act(entry *FileEntry) bool {
	strs := strings.Split(entry.base, ".")
	if len(strs) == 1 {
		entry.newpath = path.Join(entry.dir, strs[0]+a.Suffix)
	} else {
		fn := strings.Join(strs[:len(strs)-1], ".")
		if strings.HasSuffix(fn, a.Suffix) {
			return false
		}
		ext := strs[len(strs)-1]
		entry.newpath = path.Join(entry.dir, fn+a.Suffix+"."+ext)
	}
	return true
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
		format := strings.Replace(s.ReplaceFormat, "%i", strconv.Itoa(int(s.Seq.Next())), -1)
		newbase := s.Matcher.ReplaceAllString(entry.base, format)
		entry.newpath = path.Join(entry.dir, newbase)
		return true
	}
	return false
}
