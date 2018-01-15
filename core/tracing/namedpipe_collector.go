// +build linux

package tracing

import (
	"os"
	"syscall"
	"time"

	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

// record span to named pipe
func newNamedPipeCollectorLinux(path string) (zipkin.Collector, error) {
	fileInfo, err := os.Stat(path)
	needCreate := true

	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else {
		// path exists and is named pipe
		if fileInfo.Mode()&os.ModeNamedPipe == os.ModeNamedPipe {
			needCreate = false
		}
	}

	// mkfifo the file
	if needCreate {
		err = syscall.Mkfifo(path, 0640)
		if err != nil {
			return nil, err
		}
	}

	var fd *os.File
	deadLine := time.Now().Add(1 * time.Second)
	// wait the pipe reader to be ready
	for {
		fd, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE|syscall.O_NONBLOCK, os.ModeNamedPipe)
		if err == nil {
			break
		}
		if time.Now().After(deadLine) {
			break
		}
		time.Sleep(1 * time.Millisecond)
	}
	if err != nil {
		return nil, err
	}
	c := &FileCollector{
		Fd: fd,
	}
	return c, nil
}

func init() {
	newNamedPipeCollector = newNamedPipeCollectorLinux
}
