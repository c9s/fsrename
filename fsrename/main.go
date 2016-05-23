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
var fileOnlyOpt = flag.Bool("file", false, "file only")
var dirOnlyOpt = flag.Bool("dir", false, "directory only")

// replacement options
var replaceOpt = flag.String("replace", "{search}", "search")
var withOpt = flag.String("with", "{replacement}", "replacement")

/*
var extPtr = flag.String("ext", "", "extension name")
var dryRunPtr = flag.Bool("dryrun", false, "dry run only")
*/

/*
var numOfWorkersPtr = flag.Int("c", 2, "the number of concurrent rename workers. default = 2")
var trimPrefixPtr = flag.String("trimprefix", "", "trim prefix")
var trimSuffixPtr = flag.String("trimsuffix", "", "trim suffix")
var orderBy = flag.String("orderby", "", "order by")
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
		fileFilter := chain.Chain(fsrename.NewFileFilter())
		go fileFilter.Run()
		chain = fileFilter
	}
	if *dirOnlyOpt == true {
		dirFilter := chain.Chain(fsrename.NewDirFilter())
		go dirFilter.Run()
		chain = dirFilter
	}

	// string replace is enabled
	if *replaceOpt != "" {
		if *withOpt == "" {
			log.Fatalln("replacement option is required. use -with 'replacement'")
		}
		replacer := chain.Chain(fsrename.NewReplacer(*replaceOpt, *withOpt, -1))
		go replacer.Run()
		chain = replacer
	}

	printer := chain.Chain(fsrename.NewConsolePrinter())
	go printer.Run()
	chain = printer

	// send paths
	for _, path := range pathArgs {
		input <- fsrename.MustNewFileEntry(path)
	}
	input <- nil

	out := printer.Output()
	for {
		entry := <-out
		if entry == nil {
			break
		}
	}

	/*
		if *matchPatternPtr == "" {
			log.Fatalln("match pattern is required. use -match 'pattern'")
		}
		var matchRegExp = regexp.MustCompile(*matchPatternPtr)
	*/

	/*
		if *replacementPtr == "{replacement}" && *replacementFormatPtr == "{replacement}" {
			log.Fatalln("replacement is required. use -replace 'replacement' or -replace-format 'replacement with format'")
		}
		if *replacementFormatPtr != "{replacement}" {
			sequenceNumber = *seqStart
			*fileOnlyPtr = true
		}

		var extRegExp *regexp.Regexp = nil
		if *extPr != "" {
			extRegExp = regexp.MustCompile("\\." + *extPtr + "$")
		}
	*/

}
