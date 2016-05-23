package fsrename

/*
FileExtFilter is actually a regexp filter that generates the pattern from the extension name.
*/
func FileExtFilter(ext string) *RegExpFilter {
	return NewRegExpFilterWithPattern("\\." + ext + "$")
}
