// +build linux

package tracing_test

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/tracing"
	"github.com/openzipkin/zipkin-go-opentracing/thrift/gen-go/zipkincore"
	"github.com/stretchr/testify/assert"
)

func readNamedPipe(path string) {
	for {
		_, err := os.Stat(path)
		if err == nil {
			break
		}
	}
	stdout, _ := os.OpenFile(path, os.O_RDONLY, 0600)
	fmt.Println("Reading")

	bufferedReader := bufio.NewReaderSize(stdout, 16)

	byteSlice := make([]byte, 100)

	fmt.Println("Waiting for someone to write something")
	for {
		numBytesRead, err := bufferedReader.Read(byteSlice)
		if err != nil {
			return
		}
		fmt.Println("readed lenth:", numBytesRead)
		fmt.Printf("Data: %s", string(byteSlice))
	}
}

func TestNewCollector_linux(t *testing.T) {
	log.Println("Test NewCollector")

	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)

	t.Log("========new named pipe collector")
	target := "namedPipeTracing.log"
	_, err := os.Stat(target)
	if err != nil {
		assert.True(t, os.IsNotExist(err))
	} else {
		err = os.Remove(target)
		assert.NoError(t, err)
	}

	go readNamedPipe(target)
	collector, err := tracing.NewCollector(tracing.TracingNamedPipeCollector, target)
	assert.NoError(t, err)
	namedPipeCollector, ok := collector.(*tracing.FileCollector)
	assert.True(t, ok)
	assert.NotNil(t, namedPipeCollector)

	t.Log("========named pipe collector collects span")
	err = namedPipeCollector.Collect(&zipkincore.Span{Name: "test"})
	assert.NoError(t, err)
	namedPipeCollector.Close()

	t.Log("====when the named pipe file exists, not create it")
	fileInfo, err := os.Stat(target)
	assert.NoError(t, err)
	oldModifyTime := fileInfo.ModTime()

	collector, err = tracing.NewCollector(tracing.TracingNamedPipeCollector, target)
	assert.NoError(t, err)
	namedPipeCollector, ok = collector.(*tracing.FileCollector)
	assert.True(t, ok)
	assert.NotNil(t, namedPipeCollector)

	fileInfo, err = os.Stat(target)
	assert.NoError(t, err)
	newModifyTime := fileInfo.ModTime()
	assert.Equal(t, oldModifyTime, newModifyTime)

	t.Log("====write existed named pipe file")
	go readNamedPipe(target)
	err = namedPipeCollector.Collect(&zipkincore.Span{Name: "test"})
	assert.NoError(t, err)
	namedPipeCollector.Close()

	t.Log("====when write the named pipe file with no process read, err happens")
	err = namedPipeCollector.Collect(&zipkincore.Span{Name: "test"})
	assert.NotNil(t, err)
	namedPipeCollector.Close()
}
