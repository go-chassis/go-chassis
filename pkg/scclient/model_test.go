package client_test

import (
	"errors"
	"github.com/go-chassis/go-chassis/pkg/scclient"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModelException(t *testing.T) {
	t.Log("Testing modelReg.Error function")
	var modelReg *client.RegistryException = new(client.RegistryException)
	modelReg.Message = "Go-chassis"
	modelReg.Title = "fakeTitle"

	str := modelReg.Error()
	assert.Equal(t, "fakeTitle(Go-chassis)", str)

}

func TestModelExceptionOrglErr(t *testing.T) {
	t.Log("Testing modelReg.Error with title")
	var modelReg *client.RegistryException = new(client.RegistryException)
	modelReg.Message = "Go-chassis"
	modelReg.Title = "fakeTitle"
	modelReg.OrglErr = errors.New("Invalid")

	str := modelReg.Error()
	assert.Equal(t, "fakeTitle(Invalid), Go-chassis", str)

}
func TestNewCommonException(t *testing.T) {
	t.Log("Testing NewCommonException function")
	var re *client.RegistryException = new(client.RegistryException)
	re.OrglErr = nil
	re.Title = "Common exception"
	re.Message = "fakeformat"
	err := client.NewCommonException("fakeformat")
	assert.Equal(t, re, err)
}
func TestNewJsonException(t *testing.T) {
	t.Log("Testing NewJSONException function")
	var re1 *client.RegistryException = new(client.RegistryException)
	re1.OrglErr = errors.New("Invalid")
	re1.Title = "JSON exception"
	re1.Message = "args1"

	err := client.NewJSONException(errors.New("Invalid"), "args1")
	assert.Equal(t, re1, err)

	var re2 *client.RegistryException = new(client.RegistryException)
	re2.OrglErr = errors.New("Invalid")
	re2.Title = "JSON exception"
	re2.Message = ""

	err = client.NewJSONException(errors.New("Invalid"))
	assert.Equal(t, re2, err)

	var re3 *client.RegistryException = new(client.RegistryException)
	re3.OrglErr = errors.New("Invalid")
	re3.Title = "JSON exception"
	re3.Message = "[1]"

	err = client.NewJSONException(errors.New("Invalid"), 1)
	assert.Equal(t, re3, err)

}

func TestNewIOException(t *testing.T) {
	t.Log("Testing NewIOException function")
	var re1 *client.RegistryException = new(client.RegistryException)
	re1.OrglErr = errors.New("Invalid")
	re1.Title = "IO exception"
	re1.Message = "args1"

	err := client.NewIOException(errors.New("Invalid"), "args1")
	assert.Equal(t, re1, err)

	var re2 *client.RegistryException = new(client.RegistryException)
	re2.OrglErr = errors.New("Invalid")
	re2.Title = "IO exception"
	re2.Message = ""

	err = client.NewIOException(errors.New("Invalid"))
	assert.Equal(t, re2, err)

	var re3 *client.RegistryException = new(client.RegistryException)
	re3.OrglErr = errors.New("Invalid")
	re3.Title = "IO exception"
	re3.Message = "[1]"

	err = client.NewIOException(errors.New("Invalid"), 1)
	assert.Equal(t, re3, err)

}
