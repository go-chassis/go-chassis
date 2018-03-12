package main

import (
	"github.com/ServiceComb/go-chassis"
	_ "github.com/ServiceComb/go-chassis/bootstrap"
	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core"
	"github.com/ServiceComb/go-chassis/core/lager"
	"golang.org/x/net/context"
	"mime/multipart"

	"bytes"
	"io"
	"os"
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
	f, err := os.Open(filename)
	if err != nil {
		lager.Logger.Error("Error in opening file", err)
		return
	}
	defer f.Close()

	//Form part
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	fw, err := w.CreateFormFile("uploadfile", filename)
	if err != nil {
		lager.Logger.Error("Error in create form file", err)
		return
	}

	_, err = io.Copy(fw, f)
	if err != nil {
		return
	}
	w.Close()

	req, err := rest.NewRequest("POST", "cse://FileUploadServer/uploadform", buf.Bytes())

	if err != nil {
		lager.Logger.Error("new request failed.", err)
		return
	}
	defer req.Close()

	req.SetHeader("Content-Type", w.FormDataContentType())

	resp, err := core.NewRestInvoker().ContextDo(context.TODO(), req)
	if err != nil {
		lager.Logger.Error("do request failed.", err)
		return
	}
	defer resp.Close()
	lager.Logger.Info("FileUploadServer Response: " + string(resp.ReadBody()))

}
