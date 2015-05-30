package main

import "flag"
import "os"
import "fmt"
import "log"
import "path/filepath"
import "regexp"

var matchPatternPtr = flag.String("match", ".", "regular expression without slash '/'")
var replacementPtr = flag.String("replace", "", "replacement")
var fileOnlyPtr = flag.Bool("fileonly", false, "file only")
var dirOnlyPtr = flag.Bool("dironly", false, "directory only")
var forExtPtr = flag.String("forExt", "", "extension name")
var dryPtr = flag.Bool("dry", false, "dry run only")

func main() {
	flag.Parse()
	var pathArgs = flag.Args()

	// fmt.Printf("regExpPtr: %+v, fileOnlyPtr:%+v, dirOnlyPtr:%+v\n", *regExpPtr, *fileOnlyPtr, *dirOnlyPtr)
	var matchPattern = *matchPatternPtr
	var replacement = *replacementPtr
	if matchPattern == "" {
		log.Fatalln("match mattern is required. use -match 'replacement'")
	}
	if replacement == "" {
		log.Fatalln("replacement is required. use -replace 'replacement'")
	}
	var matchRegExp = regexp.MustCompile(matchPattern)

	for _, path := range pathArgs {
		var err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {

			// TODO:
			// -dironly
			// -fileonly
			// -forext

			// if it matches
			if matchRegExp.MatchString(info.Name()) {
				var newName = matchRegExp.ReplaceAllString(info.Name(), replacement)
				var newPath = filepath.Join(filepath.Dir(path), newName)
				// fmt.Printf("%s => %s\n", info.Name(), newName)
				fmt.Printf("%s => %s\n", path, newPath)
				os.Rename(path, newPath)
			}
			return err
		})
		if err != nil {
			panic(err)
		}
	}

}
