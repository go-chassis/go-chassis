package datacollect

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type ResultFile struct {
	Path   string
	fMutex sync.RWMutex
}

func (r *ResultFile) NewFile() error {
	f, err := os.Create(r.Path)
	defer f.Close()
	if err != nil {
		return err
	}
	return nil
}

func (r *ResultFile) Write(data []string) error {
	r.fMutex.Lock()
	defer r.fMutex.Unlock()
	file, err := os.OpenFile(r.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, v := range data {
		fmt.Fprintln(w, v)
	}
	return w.Flush()
}
