package workers

import "github.com/stretchr/testify/assert"
import "testing"

func TestScannerWorkerAll(t *testing.T) {
	output := make(chan *FileEntry, 10)
	cv := make(chan bool, 1)
	worker := CreateScannerWorker("tests", ScannerSettings{})
	worker(cv, nil, output)

	var entries []FileEntry
	var entry *FileEntry
	for entry = <-output; entry != nil; entry = <-output {
		entries = append(entries, *entry)
	}

	assert.Equal(t, 7, len(entries), "7 entries including all file and directory")

	ret := <-cv
	assert.True(t, ret, "cv returns true")
}

func TestScannerWorkerFileOnly(t *testing.T) {
	output := make(chan *FileEntry, 10)
	cv := make(chan bool, 1)
	worker := CreateScannerWorker("tests", ScannerSettings{FileOnly: true})
	worker(cv, nil, output)
	var entries []FileEntry
	var entry *FileEntry
	for entry = <-output; entry != nil; entry = <-output {
		entries = append(entries, *entry)
	}
	assert.Equal(t, 5, len(entries), "5 entries including all files")
	ret := <-cv
	assert.True(t, ret, "cv returns true")
}

func TestScannerWorkerDirOnly(t *testing.T) {
	output := make(chan *FileEntry, 10)
	cv := make(chan bool, 1)
	worker := CreateScannerWorker("tests", ScannerSettings{DirOnly: true})
	worker(cv, nil, output)
	var entries []FileEntry
	var entry *FileEntry
	for entry = <-output; entry != nil; entry = <-output {
		entries = append(entries, *entry)
	}
	assert.Equal(t, 2, len(entries), "2 directories")
	ret := <-cv
	assert.True(t, ret, "cv returns true")
}
