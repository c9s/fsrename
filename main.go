package main

import "flag"
import "os"
import "fmt"
import "log"
import "path/filepath"
import "regexp"
import "strings"
import "sync"
import "sort"

var matchPatternPtr = flag.String("match", ".", "regular expression without slash '/'")
var replacementPtr = flag.String("replace", "{replacement}", "replacement")
var replacementFormatPtr = flag.String("replace-format", "{replacement}", "replacement with format")
var fileOnlyPtr = flag.Bool("fileonly", false, "file only")
var dirOnlyPtr = flag.Bool("dironly", false, "directory only")
var forExtPtr = flag.String("forext", "", "extension name")
var dryRunPtr = flag.Bool("dryrun", false, "dry run only")
var numOfWorkersPtr = flag.Int("c", 2, "the number of concurrent rename workers. default = 2")
var trimPrefixPtr = flag.String("trimprefix", "", "trim prefix")
var trimSuffixPtr = flag.String("trimsuffix", "", "trim suffix")
var orderBy = flag.String("orderby", "", "order by")
var seqStart = flag.Int("seqstart", 0, "sequence number start with")
var sequenceNumber int = 1
var m sync.Mutex

type Entry struct {
	path    string
	newpath string
	info    os.FileInfo
	result  string
}

func EntryPrinter(cv chan bool, input chan *Entry) {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for {
		entry, ok := <-input
		if entry == nil || !ok {
			break
		}

		if strings.HasPrefix(entry.path, pwd) {
			var oldpath = strings.TrimLeft(strings.Replace(entry.path, pwd, "", 1), "/")
			var newpath = strings.TrimLeft(strings.Replace(entry.path, pwd, "", 1), "/")
			fmt.Printf("%s => %s", oldpath, newpath)
		} else {
			fmt.Printf("%s => %s  .. %s\n", entry.path, entry.newpath, entry.result)
		}
	}
	cv <- true
}

func GetSeqNumber() (seqNum int) {
	m.Lock()
	retValue := sequenceNumber
	sequenceNumber = sequenceNumber + 1
	m.Unlock()
	return retValue
}

func RenameWorker(cv chan bool, input chan *Entry, output chan *Entry, extRegExp *regexp.Regexp, matchRegExp *regexp.Regexp, replacement string, dryrun bool) {
	for {
		entry, ok := <-input
		if entry == nil || !ok {
			break
		}
		if extRegExp != nil && !extRegExp.MatchString(entry.info.Name()) {
			continue
		}
		if !matchRegExp.MatchString(entry.info.Name()) {
			continue
		}
		var newName string
		if *replacementPtr != "{replacement}" {
			newName = matchRegExp.ReplaceAllString(entry.info.Name(), *replacementPtr)
		} else {
			currentNumber := GetSeqNumber()
			standardFormatstring := strings.Replace(*replacementFormatPtr, "%i", "%04d", 1)
			newName = fmt.Sprintf(standardFormatstring, currentNumber)
		}

		entry.newpath = filepath.Join(filepath.Dir(entry.path), newName)
		if !dryrun {
			checkFile, err := os.Open(entry.newpath)
			defer checkFile.Close()
			if os.IsNotExist(err) {
				os.Rename(entry.path, entry.newpath)
				entry.result = " success"
			} else {
				entry.result = " file exist, ignore"
			}
		}
		output <- entry
	}
	cv <- true
}

func main() {
	flag.Parse()
	var pathArgs = flag.Args()

	if len(pathArgs) == 0 {
		pathArgs = []string{"./"}
	}

	if *replacementFormatPtr != "" {
		sequenceNumber = *seqStart
		*fileOnlyPtr = true
	}

	// Build pattern from prefix/suffix options
	if *trimPrefixPtr != "" {
		*matchPatternPtr = "^" + regexp.QuoteMeta(*trimPrefixPtr)
		*replacementFormatPtr = ""
	} else if *trimSuffixPtr != "" {
		*matchPatternPtr = regexp.QuoteMeta(*trimSuffixPtr) + "$"
	}

	if *matchPatternPtr == "" {
		log.Fatalln("match pattern is required. use -match 'pattern'")
	}
	if *replacementPtr == "{replacement}" && *replacementFormatPtr == "{replacement}" {
		log.Fatalln("replacement is required. use -replace 'replacement' or -replace-format 'replacement with format'")
	}

	var matchRegExp = regexp.MustCompile(*matchPatternPtr)
	var extRegExp *regexp.Regexp = nil
	if *forExtPtr != "" {
		extRegExp = regexp.MustCompile("\\." + *forExtPtr + "$")
	}

	var numOfWorkers = *numOfWorkersPtr

	var workerCv = make(chan bool, numOfWorkers)
	var printerCv = make(chan bool)
	var entryQueue []Entry
	var entryOutput = make(chan *Entry, 1000)
	var renamedEntryOutput = make(chan *Entry, 1000)

	for i := 0; i < numOfWorkers; i++ {
		go RenameWorker(workerCv, entryOutput, renamedEntryOutput, extRegExp, matchRegExp, *replacementPtr, *dryRunPtr)
	}
	go EntryPrinter(printerCv, renamedEntryOutput)

	for _, pathArg := range pathArgs {
		matches, err := filepath.Glob(pathArg)
		if err != nil {
			panic(err)
		}

		for _, match := range matches {
			var err = filepath.Walk(match, func(path string, info os.FileInfo, err error) error {
				if *dirOnlyPtr {
					if info.IsDir() {
						entryQueue = append(entryQueue, Entry{path: path, info: info})
					}
				} else if *fileOnlyPtr {
					if !info.IsDir() {
						entryQueue = append(entryQueue, Entry{path: path, info: info})
					}
				} else {
					entryQueue = append(entryQueue, Entry{path: path, info: info})
				}
				return err
			})
			if err != nil {
				panic(err)
			}
		}

		//Sorting only enable on file only. Default sorting is alphabet
		if *fileOnlyPtr && *orderBy != "" {
			switch *orderBy {
			case "Reverse":
				sort.Sort(ReverseSort{entryQueue})
			case "Mtime":
				sort.Sort(MtimeSort{entryQueue})
			case "MtimeReverse":
				sort.Sort(MtimeReverseSort{entryQueue})
			case "Size":
				sort.Sort(SizeSort{entryQueue})
			case "SizeReverse":
				sort.Sort(SizeReverseSort{entryQueue})
			}
		}

		for index, _ := range entryQueue {
			entryOutput <- &(entryQueue[index])
		}
	}
	entryOutput <- nil
	close(entryOutput)

	for ; numOfWorkers > 0; numOfWorkers-- {
		<-workerCv
	}
	close(workerCv)
	renamedEntryOutput <- nil

	<-printerCv
	close(printerCv)
}
