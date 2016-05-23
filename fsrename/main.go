package main

import "flag"
import "os"
import "fmt"
import "log"
import "regexp"
import "github.com/c9s/fsrename"

// filter options
var matchOpt = flag.String("match", "", "pre-filter (regular expression without slash '/')")
var containsOpt = flag.String("contains", "", "strings.contains filter")
var extOpt = flag.String("ext", "", "extension name filter")
var fileOnlyOpt = flag.Bool("file", false, "file only")
var dirOnlyOpt = flag.Bool("dir", false, "directory only")

// replacement options
var replaceOpt = flag.String("replace", "{nil}", "search")
var rOpt = flag.String("r", "{nil}", "search")
var replaceRegexpOpt = flag.String("replaceRegexp", "{nil}", "regular expression replace target")
var rreOpt = flag.String("rre", "{nil}", "regular expression replace target")

var withOpt = flag.String("with", "{nil}", "replacement")
var wOpt = flag.String("w", "{nil}", "replacement")
var withFormatOpt = flag.String("withFormat", "{nil}", "replacement format")

// rule builders
var trimPrefixOpt = flag.String("trimPrefix", "", "trim prefix")
var trimSuffixOpt = flag.String("trimSuffix", "", "trim suffix")
var camelOpt = flag.Bool("camel", false, "Convert substrings to camel cases")
var underscoreOpt = flag.Bool("underscore", false, "Convert substrings to underscore cases")

// runtime option
var dryRunOpt = flag.Bool("dryrun", false, "dry run only")
var orderOpt = flag.String("order", "", "order by")

/*
var seqStart = flag.Int("seqstart", 0, "sequence number start with")
*/
func main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	flag.Parse()
	var pathArgs = flag.Args()

	// runs without any arguments, find files under the current directory
	if len(pathArgs) == 0 {
		pathArgs = []string{pwd}
	}

	var chain fsrename.Worker

	input := fsrename.NewFileStream()
	scanner := fsrename.NewGlobScanner()
	scanner.SetInput(input)
	chain = scanner
	go scanner.Run()

	if *fileOnlyOpt == true {
		chain = chain.Chain(fsrename.NewFileFilter())
		go chain.Run()
	}
	if *dirOnlyOpt == true {
		chain = chain.Chain(fsrename.NewDirFilter())
		go chain.Run()
	}
	if *extOpt != "" {
		chain = chain.Chain(fsrename.NewFileExtFilter(*extOpt))
		go chain.Run()
	}

	if *matchOpt != "" {
		chain = chain.Chain(fsrename.NewRegExpFilterWithPattern(*matchOpt))
		go chain.Run()
	}
	if *containsOpt != "" {
		chain = chain.Chain(fsrename.NewStrContainsFilter(*containsOpt))
		go chain.Run()
	}

	if *fileOnlyOpt && *orderOpt != "" {
		switch *orderOpt {
		case "reverse":
			chain = chain.Chain(fsrename.NewReverseSorter())
			go chain.Run()
			break
		case "mtime":
			chain = chain.Chain(fsrename.NewMtimeSorter())
			go chain.Run()
			break
		case "reverse-mtime":
			chain = chain.Chain(fsrename.NewMtimeReverseSorter())
			go chain.Run()
			break
		case "size":
			chain = chain.Chain(fsrename.NewSizeSorter())
			go chain.Run()
			break
		case "reverse-size":
			chain = chain.Chain(fsrename.NewSizeReverseSorter())
			go chain.Run()
			break
		}
	}

	// copy short option to long option
	if *rOpt != "{nil}" {
		*replaceOpt = *rOpt
	}
	if *rreOpt != "{nil}" {
		*replaceRegexpOpt = *rreOpt
	}
	if *wOpt != "{nil}" {
		*withOpt = *wOpt
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
		go chain.Run()
	} else if *underscoreOpt == true {
		chain = chain.Chain(fsrename.NewUnderscoreReplacer())
		go chain.Run()
	}

	// string replace is enabled
	if *replaceOpt != "{nil}" || *replaceRegexpOpt != "{nil}" {
		if *withOpt == "{nil}" && *withFormatOpt == "{nil}" {
			log.Fatalln("replacement option is required. use -with 'replacement' or -withFormat 'format'.")
		}

		if *replaceRegexpOpt != "{nil}" {
			if *withOpt != "{nil}" {
				chain = chain.Chain(fsrename.NewRegExpReplacer(*replaceRegexpOpt, *withOpt))
				go chain.Run()
			} else if *withFormatOpt != "{nil}" {
				chain = chain.Chain(fsrename.NewRegExpFormatReplacer(*replaceRegexpOpt, *withFormatOpt))
				go chain.Run()
			}
		} else {
			if *withOpt != "{nil}" {
				chain = chain.Chain(fsrename.NewStrReplacer(*replaceOpt, *withOpt, -1))
				go chain.Run()
			} else if *withFormatOpt != "{nil}" {
				chain = chain.Chain(fsrename.NewFormatReplacer(*replaceOpt, *withFormatOpt))
				go chain.Run()
			}
		}

	}

	if *dryRunOpt == false {
		chain = chain.Chain(fsrename.NewRenamer())
		go chain.Run()
	}

	chain = chain.Chain(fsrename.NewConsolePrinter())
	go chain.Run()

	// send paths
	for _, path := range pathArgs {
		input <- fsrename.MustNewFileEntry(path)
	}
	input <- nil

	out := chain.Output()
	for {
		entry := <-out
		if entry == nil {
			break
		}
	}
}
