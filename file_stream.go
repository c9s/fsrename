package fsrename

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

// A channel of file entries
type FileStream chan *FileEntry

// Create a channel for sending file entries
func NewFileStream() FileStream {
	return make(FileStream, 10)
}

func FileStreamFromChangeLog(input FileStream, file string) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	for {
		columns, err := r.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}

		// columns: entry.path, entry.newpath, entry.message
		entry, err := NewFileEntry(columns[1])
		if err != nil {
			log.Println(columns[1], err)
			// get next row
			continue
		}
		// reverse the rename
		entry.newpath = columns[0]
		input <- entry
	}
	// end the list
	input <- nil
}
