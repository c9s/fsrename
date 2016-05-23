package workers

/*
SuffixFilter is actually a regexp filter that generates the pattern from the prefix
*/
func SuffixFilter(suffix string) *RegExpFilter {
	return NewRegExpFilterWithPattern(suffix + "$")
}
