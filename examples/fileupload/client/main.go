package main

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"os"

	"fmt"
	"github.com/ServiceComb/go-chassis"
	_ "github.com/ServiceComb/go-chassis/bootstrap"
	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core"
	"github.com/ServiceComb/go-chassis/core/lager"
	"io/ioutil"
)

func main() {
	//Init framework
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed.", err)
		return
	}

	// file / form to upload
	uploadfile("file.input")
	uploadform("form.input")
}

func uploadfile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		lager.Logger.Error("Error in opening file", err)
		return
	}
	defer f.Close()

	body := &bytes.Buffer{}

	_, err = io.Copy(body, f)
	if err != nil {
		lager.Logger.Error("Copy failed.", err)
		return
	}

	req, err := rest.NewRequest("POST", "cse://FileUploadServer/uploadfile", body.Bytes())
	if err != nil {
		lager.Logger.Error("new request failed.", err)
		return
	}
	defer req.Close()

	req.SetHeader("Content-Type", "application/octet-stream")

	resp, err := core.NewRestInvoker().ContextDo(context.TODO(), req)
	if err != nil {
		lager.Logger.Error("do request failed.", err)
		return
	}
	defer resp.Close()
	lager.Logger.Info("FileUploadServer Response: " + string(resp.ReadBody()))

}

func uploadform(filename string) {
	//Form part
	headBuf := bytes.NewBufferString("")
	headBufWriter := multipart.NewWriter(headBuf)
	_, err := headBufWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		lager.Logger.Error("Error in create form file", err)
		return
	}

	f, err := os.Open(filename)
	if err != nil {
		lager.Logger.Error("Error in opening file", err)
		return
	}
	defer f.Close()

	fs, err := f.Stat()
	if err != nil {
		lager.Logger.Error("Error in stat file", err)
		return
	}

	lastBoundary := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", headBufWriter.Boundary()))

	bodyReader := io.MultiReader(headBuf, f, lastBoundary)

	req, err := rest.NewRequest("POST", "cse://FileUploadServer/uploadform")
	if err != nil {
		lager.Logger.Error("new request failed.", err)
		return
	}
	req.Req.Body = ioutil.NopCloser(bodyReader)
	req.SetHeader("Content-Type", headBufWriter.FormDataContentType())
	req.Req.ContentLength = int64(headBuf.Len()) + fs.Size() + int64(lastBoundary.Len())

	defer req.Close()

	resp, err := core.NewRestInvoker().ContextDo(context.TODO(), req)
	if err != nil {
		lager.Logger.Error("do request failed.", err)
		return
	}
	defer resp.Close()
	lager.Logger.Info("FileUploadServer Response: " + string(resp.ReadBody()))

}
