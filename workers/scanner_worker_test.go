package workers

import "github.com/stretchr/testify/assert"
import "testing"

func TestScannerWorker(t *testing.T) {
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
