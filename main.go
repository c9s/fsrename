package main

import "flag"
import "os"
import "fmt"
import "log"
import "path/filepath"
import "regexp"
import "strings"

var matchPatternPtr = flag.String("match", ".", "regular expression without slash '/'")
var replacementPtr = flag.String("replace", "", "replacement")
var fileOnlyPtr = flag.Bool("fileonly", false, "file only")
var dirOnlyPtr = flag.Bool("dironly", false, "directory only")
var forExtPtr = flag.String("forExt", "", "extension name")
var dryPtr = flag.Bool("dry", false, "dry run only")

type Entry struct {
	path    string
	newpath string
	info    os.FileInfo
}

func EntryPrinter(cv chan bool, input chan *Entry) {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// fmt.Println(pwd)

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
			fmt.Printf("%s => %s\n", entry.path, entry.newpath)
		}
	}
	cv <- true
}

func RenameWorker(cv chan bool, input chan *Entry, output chan *Entry, matchRegExp *regexp.Regexp, replacement string) {
	for {
		entry, ok := <-input
		if entry == nil || !ok {
			break
		}
		var newName = matchRegExp.ReplaceAllString(entry.info.Name(), *replacementPtr)
		entry.newpath = filepath.Join(filepath.Dir(entry.path), newName)
		os.Rename(entry.path, entry.newpath)
		output <- entry
	}
	cv <- true
}

func main() {
	flag.Parse()
	var pathArgs = flag.Args()

	if *matchPatternPtr == "" {
		log.Fatalln("match mattern is required. use -match 'pattern'")
	}
	if *replacementPtr == "" {
		log.Fatalln("replacement is required. use -replace 'replacement'")
	}
	var matchRegExp = regexp.MustCompile(*matchPatternPtr)

	var numOfWorkers = 3

	var workerCv = make(chan bool, numOfWorkers)
	var printerCv = make(chan bool)
	var entryOutput = make(chan *Entry, 1000)
	var renamedEntryOutput = make(chan *Entry, 1000)

	for i := 0; i < numOfWorkers; i++ {
		go RenameWorker(workerCv, entryOutput, renamedEntryOutput, matchRegExp, *replacementPtr)
	}
	go EntryPrinter(printerCv, renamedEntryOutput)

	for _, path := range pathArgs {
		var err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			// TODO:
			// -dironly
			// -fileonly
			// -forext
			if matchRegExp.MatchString(info.Name()) {
				entryOutput <- &Entry{path: path, info: info}
			}
			return err
		})
		if err != nil {
			panic(err)
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
