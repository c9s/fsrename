package fsrename

import "regexp"

func NewStrReplacer(search, replace string, n int) *Actor {
	return NewActor(NewStrReplaceAction(search, replace, n))
}

func NewFormatReplacer(search, replaceFormat string) *Actor {
	return NewActor(NewStrFormatReplaceAction(search, replaceFormat))
}

func NewCamelCaseReplacer() *Actor {
	return NewActor(NewCamelCaseAction("[-_]+"))
}

func NewUnderscoreReplacer() *Actor {
	return NewActor(NewUnderscoreAction())
}

func NewPrefixAdder(prefix string) *Actor {
	return NewActor(NewPrefixAction(prefix))
}

func NewSuffixAdder(suffix string) *Actor {
	return NewActor(NewSuffixAction(suffix))
}

func NewRegExpReplacer(pattern, replace string) *Actor {
	matcher := regexp.MustCompile(pattern)
	return NewActor(NewRegExpReplaceAction(matcher, replace))
}

func NewRegExpFormatReplacer(pattern, replace string) *Actor {
	matcher := regexp.MustCompile(pattern)
	return NewActor(NewRegExpFormatReplaceAction(matcher, replace))
}
