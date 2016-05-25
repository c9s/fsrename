// Copyright 2014-2016 Yo-An Lin. All rights reserved.
// license that can be found in the LICENSE file.

/*
NAME

	fsrename

SYNOPSIS

    fsrename [options] [path...]

DESCRIPTION

	When `[path...]` is not given, "./" will be used as the default path for scanning files.

OPTIONS

	FILTER OPTIONS

	-match

		pre-filter the files and directories based on the given regular
		pattern.

	-contains

		pre-filter the files and directories based on the given string needle.

	-file, -f

		only for files.

	-dir, -d

		only for directories.

	-ext

		find files with matched file extension.


	REPLACEMENT OPTIONS

	Please note the replacement target only works for the basename of a path.
	-replace* and -with* should be combined together to replace the substrings.

		-replace, -r

			specify target substring with normal string matching.

		-replace-regexp, -rr

			specify target substring with regular expression matching.

		-with, -w

			replacement for the target substring.

		-with-format

			replacement with fmt.Sprintf format for the target substring.

	REPLACE RULE BUILDER OPTIONS

	-trim-prefix

		trim filename prefix.

	-trim-suffix

		trim filename suffix (this option removes suffix even for filename
		extensions).

	-camel

		converts dash/underscore separated filenames into camelcase filenames.

	-underscore

		converts camelcase filesnames into underscore separated filenames.

	COMMON OPTIONS

		-dryrun

			dry run, don't rename, just preview the result.

	QUICK EXAMPLES

	Find files with extension `.php` and replace the substring from the filename.

		fsrename -ext "php" -replace "some" -with "others" src/

	Replace `Stmt.go` with "_stmt.go" under the current directory:

		fsrename -replace "Stmt.go" -with "_stmt.go"

	Replace `Stmt.go` with "_stmt.go" under directory `src/c6`:

		fsrename -replace "Stmt.go" -with "_stmt.go" src/c6

	Replace `foo` with `bar` from files contains `prefix_`

		fsrename -file -contains prefix_ -replace foo -with bar test

	Or use `-match` to pre-filter the files with regular expression

		fsrename -file -match "[a-z]+" -replace foo -with bar test

	Use regular expression without escaping:

		fsrename -replace-regexp "_[a-z]*.go" -with ".go" src/c6

		fsrename -file -replace-regexp "_[a-z]*.go" -with ".go" src/c6

		fsrename -file -ext go -replace-regexp "[a-z]*" -with "123" src/c6

		fsrename -dir -replace "_xxx" -with "_aaa" src/c6

		fsrename -replace "_xxx" -with "_aaa" -dryrun  src/c6
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
)

import "github.com/c9s/fsrename"

// filter options
var matchOpt = flag.String("match", "", "pre-filter (regular expression without slash '/')")
var containsOpt = flag.String("contains", "", "strings.contains filter")
var extOpt = flag.String("ext", "", "extension name filter")
var fileOnlyOpt = flag.Bool("file", false, "file only")
var fOpt = flag.Bool("f", false, "an alias of file only")
var dirOnlyOpt = flag.Bool("dir", false, "directory only")
var dOpt = flag.Bool("d", false, "an alias of dir only")

var changelogOpt = flag.String("changelog", "", "the changelog file")

// replacement options
var replaceOpt = flag.String("replace", "{nil}", "search")
var rOpt = flag.String("r", "{nil}", "search")
var replaceRegexpOpt = flag.String("replace-regexp", "{nil}", "regular expression replace target")
var replaceReOpt = flag.String("replace-re", "{nil}", "regular expression replace target")
var rrOpt = flag.String("rr", "{nil}", "regular expression replace target")

var withOpt = flag.String("with", "{nil}", "replacement")
var wOpt = flag.String("w", "{nil}", "replacement")
var withFormatOpt = flag.String("with-format", "{nil}", "replacement format")
var wfOpt = flag.String("wf", "{nil}", "replacement format")

// rule builders
var trimPrefixOpt = flag.String("trim-prefix", "", "trim prefix")
var trimSuffixOpt = flag.String("trim-suffix", "", "trim suffix")
var camelOpt = flag.Bool("camel", false, "Convert substrings to camel cases")
var underscoreOpt = flag.Bool("underscore", false, "Convert substrings to underscore cases")

var rollbackOpt = flag.String("rollback", "", "rollback renames from a changelog file")

// runtime option
var dryRunOpt = flag.Bool("dryrun", false, "dry run only")

var orderOpt = flag.String("order", "", "order by")

/*
var seqStart = flag.Int("seqstart", 0, "sequence number start with")
*/
func main() {
	flag.Parse()

	var chain fsrename.Worker

	input := fsrename.NewFileStream()

	if *rollbackOpt != "" {
		chain = fsrename.NewProxy()
	} else {
		chain = fsrename.NewGlobScanner()
	}
	chain.SetInput(input)
	chain.Start()

	// copy short option to long option
	if *rOpt != "{nil}" {
		*replaceOpt = *rOpt
	}
	if *rrOpt != "{nil}" {
		*replaceRegexpOpt = *rrOpt
	} else if *replaceReOpt != "{nil}" {
		*replaceRegexpOpt = *replaceReOpt
	}

	if *wOpt != "{nil}" {
		*withOpt = *wOpt
	}
	if *wfOpt != "{nil}" {
		*withFormatOpt = *wfOpt
	}

	if *fOpt == true {
		*fileOnlyOpt = true
	}
	if *dOpt == true {
		*dirOnlyOpt = true
	}

	if *fileOnlyOpt == true {
		chain = chain.Chain(fsrename.NewFileFilter())
	}
	if *dirOnlyOpt == true {
		chain = chain.Chain(fsrename.NewDirFilter())
	}
	if *extOpt != "" {
		chain = chain.Chain(fsrename.NewFileExtFilter(*extOpt))
	}

	if *matchOpt != "" {
		chain = chain.Chain(fsrename.NewRegExpFilterWithPattern(*matchOpt))
	}
	if *containsOpt != "" {
		chain = chain.Chain(fsrename.NewStrContainsFilter(*containsOpt))
	}

	if *fileOnlyOpt && *orderOpt != "" {
		switch *orderOpt {
		case "reverse":
			chain = chain.Chain(fsrename.NewReverseSorter())
			break
		case "mtime":
			chain = chain.Chain(fsrename.NewMtimeSorter())
			break
		case "reverse-mtime":
			chain = chain.Chain(fsrename.NewMtimeReverseSorter())
			break
		case "size":
			chain = chain.Chain(fsrename.NewSizeSorter())
			break
		case "reverse-size":
			chain = chain.Chain(fsrename.NewSizeReverseSorter())
			break
		}
	}

	if *trimPrefixOpt != "" {
		*replaceRegexpOpt = "^" + regexp.QuoteMeta(*trimPrefixOpt)
		*withOpt = ""
	}
	if *trimSuffixOpt != "" {
		*replaceRegexpOpt = regexp.QuoteMeta(*trimSuffixOpt) + "$"
		*withOpt = ""
	}

	if *camelOpt == true {
		chain = chain.Chain(fsrename.NewCamelCaseReplacer())
	} else if *underscoreOpt == true {
		chain = chain.Chain(fsrename.NewUnderscoreReplacer())
	}

	// string replace is enabled
	if *replaceOpt != "{nil}" || *replaceRegexpOpt != "{nil}" {
		if *withOpt == "{nil}" && *withFormatOpt == "{nil}" {
			log.Fatalln("replacement option is required. use -with 'replacement' or -with-format 'format'.")
		}

		if *replaceRegexpOpt != "{nil}" {
			if *withOpt != "{nil}" {
				chain = chain.Chain(fsrename.NewRegExpReplacer(*replaceRegexpOpt, *withOpt))
			} else if *withFormatOpt != "{nil}" {
				chain = chain.Chain(fsrename.NewRegExpFormatReplacer(*replaceRegexpOpt, *withFormatOpt))
			}
		} else {
			if *withOpt != "{nil}" {
				chain = chain.Chain(fsrename.NewStrReplacer(*replaceOpt, *withOpt, -1))
			} else if *withFormatOpt != "{nil}" {
				chain = chain.Chain(fsrename.NewFormatReplacer(*replaceOpt, *withFormatOpt))
			}
		}

	}

	// Always run renamer at the end
	chain = chain.Chain(fsrename.NewRenamer(*dryRunOpt))

	if *changelogOpt != "" {
		chain = chain.Chain(fsrename.NewChangeLogWriter(*changelogOpt))
	}

	chain = chain.Chain(fsrename.NewConsolePrinter())

	if *rollbackOpt != "" {
		// send file rename entries from csv files
		fsrename.FileStreamFromChangeLog(input, *rollbackOpt)
	} else {
		var pathArgs = flag.Args()
		if len(pathArgs) == 0 {
			// runs without any arguments, find files under the current directory
			pwd, err := os.Getwd()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			pathArgs = []string{pwd}
		}
		// send paths
		for _, path := range pathArgs {
			input <- fsrename.MustNewFileEntry(path)
		}
		input <- nil
	}

	// TODO: use condvar instead receiving the paths...
	out := chain.Output()
	for {
		entry := <-out
		if entry == nil {
			break
		}
	}
}
