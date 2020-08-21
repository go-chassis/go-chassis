package lager

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyFile(t *testing.T) {
	err := ioutil.WriteFile("./test.copy", []byte("test"), 0700)
	assert.NoError(t, err)
	err = CopyFile("./test.copy", "./test2.copy")
	assert.NoError(t, err)
	b, err := ioutil.ReadFile("./test2.copy")
	assert.NoError(t, err)
	assert.Equal(t, []byte("test"), b)
}

func Test_removeFile(t *testing.T) {
	ioutil.WriteFile("./remove.copy", []byte("test"), 0700)
	p := "remove.copy"

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"remove", args{p}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := removeFile(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("removeFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
