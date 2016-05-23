package main

import "flag"
import "os"
import "fmt"
import "log"

// import "path/filepath"
// import "sync"
import "github.com/c9s/fsrename"

// filter options
var matchOpt = flag.String("match", ".", "regular expression without slash '/'")
var extOpt = flag.String("ext", "", "extension name filter")
var fileOnlyOpt = flag.Bool("file", false, "file only")
var dirOnlyOpt = flag.Bool("dir", false, "directory only")

// replacement options
var replaceOpt = flag.String("replace", "", "search")
var replaceRegexpOpt = flag.String("replaceRegexp", "", "search")
var withOpt = flag.String("with", "", "replacement")
var withFormatOpt = flag.String("withFormat", "", "replacement format")

// rule builders
var trimPrefixOpt = flag.String("trimprefix", "", "trim prefix")
var trimSuffixOpt = flag.String("trimsuffix", "", "trim suffix")

// runtime option
var dryRunOpt = flag.Bool("dryrun", false, "dry run only")
var orderOpt = flag.String("order", "", "order by")

/*
var numOfWorkersPtr = flag.Int("c", 2, "the number of concurrent rename workers. default = 2")
var seqStart = flag.Int("seqstart", 0, "sequence number start with")
var sequenceNumber int = 1
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

	// string replace is enabled
	if *replaceOpt != "" || *replaceRegexpOpt != "" {
		if *withOpt == "" && *withFormatOpt == "" {
			log.Fatalln("replacement option is required. use -with 'replacement' or -withFormat 'format'.")
		}

		if *replaceRegexpOpt != "" {
			if *withOpt != "" {
				chain = chain.Chain(fsrename.NewRegExpReplacer(*replaceRegexpOpt, *withOpt))
				go chain.Run()
			} else if *withFormatOpt != "" {
				chain = chain.Chain(fsrename.NewRegExpFormatReplacer(*replaceRegexpOpt, *withFormatOpt))
				go chain.Run()
			}
		} else {
			if *withOpt != "" {
				chain = chain.Chain(fsrename.NewStrReplacer(*replaceOpt, *withOpt, -1))
				go chain.Run()
			} else if *withFormatOpt != "" {
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
