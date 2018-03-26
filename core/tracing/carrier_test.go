package tracing_test

import (
	"testing"

	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core/tracing"
	"github.com/stretchr/testify/assert"
)

func TestRestClientHeaderWriter(t *testing.T) {
	key := "key1"
	val := "val1"
	restClientHeader, err := rest.NewRequest("GET", "localhost", make([]byte, 0))
	assert.NoError(t, err)
	restClientHeader.SetHeader(key, val)
	assert.Equal(t, val, restClientHeader.GetHeader(key))

	restClientHeaderWriter := (*tracing.RestClientHeaderWriter)(restClientHeader)
	//txtMapWriter := opentracing.TextMapWriter(restClientHeaderWriter)
	//txtMapWriter.Set(key, val)
	restClientHeaderWriter.Set(key, val)
	assert.Equal(t, val, restClientHeader.GetHeader(key))
}

func TestHeaderCarrier(t *testing.T) {
	key := "x-b3-traceid"
	val := "abc"
	headerMap := make(map[string]string)
	headerMap[key] = val
	carrier := &tracing.HeaderCarrier{Header: headerMap}

	containsZipkinHeader := false
	handlerFunc := func(k, v string) error {
		if k == key || v == val {
			containsZipkinHeader = true
			return nil
		}
		return nil
	}
	carrier.ForeachKey(handlerFunc)
	assert.Equal(t, true, containsZipkinHeader)
}
