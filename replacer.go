package fsrename

import "regexp"

func NewStrReplacer(search, replace string, n int) *Actor {
	return NewActor(NewStrReplaceAction(search, replace, n))
}

func NewFormatReplacer(search, replaceFormat string) *Actor {
	return NewActor(NewStrFormatReplaceAction(search, replaceFormat))
}

func NewRegExpReplacer(pattern, replace string) *Actor {
	matcher := regexp.MustCompile(pattern)
	return NewActor(NewRegExpReplaceAction(matcher, replace))
}

func NewRegExpFormatReplacer(pattern, replace string) *Actor {
	matcher := regexp.MustCompile(pattern)
	return NewActor(NewRegExpFormatReplaceAction(matcher, replace))
}
