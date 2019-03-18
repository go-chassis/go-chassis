package lager_test

import (
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestCopyFile(t *testing.T) {
	err := ioutil.WriteFile("./test.copy", []byte("test"), 0600)
	assert.NoError(t, err)
	err = lager.CopyFile("./test.copy", "./test2.copy")
	assert.NoError(t, err)
	b, err := ioutil.ReadFile("./test2.copy")
	assert.NoError(t, err)
	assert.Equal(t, []byte("test"), b)
}
