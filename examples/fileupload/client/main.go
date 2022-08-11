package main

import (
	"bytes"
	"context"
	"github.com/go-chassis/openlog"
	"io"
	"mime/multipart"
	"os"

	"fmt"
	"github.com/go-chassis/go-chassis/v2"
	_ "github.com/go-chassis/go-chassis/v2/bootstrap"
	"github.com/go-chassis/go-chassis/v2/client/rest"
	"github.com/go-chassis/go-chassis/v2/core"
	"github.com/go-chassis/go-chassis/v2/pkg/util/httputil"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/fileupload/client/
func main() {
	//Init framework
	if err := chassis.Init(); err != nil {
		openlog.Error("Init failed." + err.Error())
		return
	}

	// file / form to upload
	uploadfile("file.input")
	uploadform("form.input")
}

func uploadfile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		openlog.Error("Error in opening file" + err.Error())
		return
	}
	defer f.Close()

	body := &bytes.Buffer{}

	_, err = io.Copy(body, f)
	if err != nil {
		openlog.Error("Copy failed." + err.Error())
		return
	}

	req, err := rest.NewRequest("POST", "http://FileUploadServer/uploadfile", body.Bytes())
	if err != nil {
		openlog.Error("new request failed." + err.Error())
		return
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := core.NewRestInvoker().ContextDo(context.TODO(), req)
	if err != nil {
		openlog.Error("do request failed." + err.Error())
		return
	}
	defer resp.Body.Close()
	openlog.Info("FileUploadServer Response: " + string(httputil.ReadBody(resp)))

}

func uploadform(filename string) {
	//Form part
	headBuf := bytes.NewBufferString("")
	headBufWriter := multipart.NewWriter(headBuf)
	_, err := headBufWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		openlog.Error("Error in create form file" + err.Error())
		return
	}

	f, err := os.Open(filename)
	if err != nil {
		openlog.Error("Error in opening file" + err.Error())
		return
	}
	defer f.Close()

	fs, err := f.Stat()
	if err != nil {
		openlog.Error("Error in stat file" + err.Error())
		return
	}

	lastBoundary := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", headBufWriter.Boundary()))

	bodyReader := io.MultiReader(headBuf, f, lastBoundary)

	req, err := rest.NewRequest("POST", "http://FileUploadServer/uploadform", nil)
	if err != nil {
		openlog.Error("new request failed." + err.Error())
		return
	}
	req.Body = io.NopCloser(bodyReader)
	req.Header.Set("Content-Type", headBufWriter.FormDataContentType())
	req.ContentLength = int64(headBuf.Len()) + fs.Size() + int64(lastBoundary.Len())

	resp, err := core.NewRestInvoker().ContextDo(context.TODO(), req)
	if err != nil {
		openlog.Error("do request failed." + err.Error())
		return
	}
	defer resp.Body.Close()
	openlog.Info("FileUploadServer Response: " + string(httputil.ReadBody(resp)))

}
