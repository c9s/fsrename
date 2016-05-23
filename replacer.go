package fsrename

func NewReplacer(search, replace string, n int) *Actor {
	return NewActor(NewStrReplaceAction(search, replace, n))
}

func NewFormatReplacer(search, replaceFormat string) *Actor {
	return NewActor(NewStrFormatReplaceAction(search, replaceFormat))
}

func NewRegExpReplacer(pattern, replace string) *Actor {
	return NewActor(NewRegExpActionWithPattern(pattern, replace))
}
