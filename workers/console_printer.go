package workers

import "fmt"
import "strings"
import "os"

type ConsolePrinter struct {
	*BaseWorker
}

func NewConsolePrinter() *ConsolePrinter {
	return &ConsolePrinter{NewBaseWorker()}
}

func (w *ConsolePrinter) Run() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for {
		select {
		case <-w.stop:
			return
		case entry := <-w.input:
			// end of data
			if entry == nil {
				w.emitEnd()
				return
			}
			// trim pwd paths
			if strings.HasPrefix(entry.path, pwd) {
				var oldpath = strings.TrimLeft(strings.Replace(entry.path, pwd, "", 1), "/")
				var newpath = strings.TrimLeft(strings.Replace(entry.path, pwd, "", 1), "/")
				fmt.Printf("'%s' => '%s'\n", oldpath, newpath)
			} else {
				fmt.Printf("'%s' => '%s'\n", entry.path, entry.newpath)
			}
		}
	}
}
