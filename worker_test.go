package fsrename

import "testing"
import "github.com/stretchr/testify/assert"

func TestScanner(t *testing.T) {
	input := NewFileStream()
	output := NewFileStream()
	worker := NewGlobScanner()
	worker.SetInput(input)
	worker.SetOutput(output)
	go worker.Run()
	input <- &FileEntry{"tests/scanner", nil, "", nil}
	input <- nil
	assert.NotNil(t, output)

	entries := readEntries(t, output)
	assert.Equal(t, 7, len(entries))
}

func TestFilterWorkerEmtpy(t *testing.T) {
	input := NewFileStream()
	output := NewFileStream()
	worker := NewFileFilter()
	worker.SetInput(input)
	worker.SetOutput(output)
	go worker.Run()

	e, err := NewFileEntry("tests/autoload.php")
	assert.Nil(t, err)
	input <- e
	input <- nil
	assert.NotNil(t, output)
	entries := readEntries(t, output)
	assert.Equal(t, 1, len(entries))
}

func TestSimpleRegExpPipe(t *testing.T) {
	scanner := NewGlobScanner()
	filter := scanner.Chain(NewFileFilter())
	filter2 := filter.Chain(NewRegExpFilterWithPattern("\\.php$"))

	go scanner.Run()
	go filter.Run()
	go filter2.Run()

	input := NewFileStream()
	scanner.SetInput(input)
	input <- &FileEntry{"tests", nil, "", nil}
	input <- nil
	output := filter2.Output()
	assert.NotNil(t, output)
	entries := readEntries(t, output)
	assert.Equal(t, 2, len(entries))
}

func TestSimpleFilePipe(t *testing.T) {
	scanner := NewGlobScanner()
	filter := scanner.Chain(&FileFilter{
		&BaseWorker{stop: make(CondVar, 1)},
	})

	go scanner.Run()
	go filter.Run()

	input := NewFileStream()
	scanner.SetInput(input)

	input <- MustNewFileEntry("tests/scanner")
	input <- nil
	output := filter.Output()
	assert.NotNil(t, output)
	entries := readEntries(t, output)
	assert.Equal(t, 5, len(entries))
}

func TestSimpleReverseSorter(t *testing.T) {
	scanner := NewGlobScanner()
	filter := scanner.Chain(&FileFilter{NewBaseWorker()})
	sorter := filter.Chain(&ReverseSorter{NewBaseWorker()})
	go scanner.Run()
	go filter.Run()
	go sorter.Run()

	input := NewFileStream()
	scanner.SetInput(input)

	input <- MustNewFileEntry("tests")
	input <- nil

	output := sorter.Output()
	assert.NotNil(t, output)
	entries := readEntries(t, output)
	assert.Equal(t, 10, len(entries))
}

func readEntries(t *testing.T, input chan *FileEntry) []FileEntry {
	var entries []FileEntry
	for {
		select {
		case entry := <-input:
			t.Log("entry", entry)
			if entry == nil {
				return entries
			}
			entries = append(entries, *entry)
			break
		}
	}
	return entries
}
