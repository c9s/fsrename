package fsrename

import "regexp"

/*
PrefixFilter is actually a regexp filter that generates the pattern from the prefix
*/
func PrefixFilter(prefix string) *RegExpFilter {
	return NewRegExpFilterWithPattern("^" + regexp.QuoteMeta(prefix))
}
