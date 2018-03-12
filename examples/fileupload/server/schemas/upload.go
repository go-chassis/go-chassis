package schemas

import (
	"errors"
	rf "github.com/ServiceComb/go-chassis/server/restful"
	"io"
	"net/http"
	"os"
)

//RestFulUpload is a struct used to implement restful upload
type RestFulUpload struct {
}

// UploadFile is a method used to reply user hello in json format
func (r *RestFulUpload) UploadFile(b *rf.Context) {
	//USAGE: curl -X POST http://127.0.0.1:8083/uploadfile -H 'content-type: application/octet-stream' --data-binary '@input.txt'
	req := b.ReadRequest()
	res := b.ReadResponseWriter()
	err := uploadFile(res, req)
	if err == nil {
		b.Write([]byte("Upload file OK!"))
	}
}

func uploadFile(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	//POST takes the uploaded file(s) and saves it to disk.
	case http.MethodPost:
		var count int64
		var ntest int64
		defer r.Body.Close()
		buf := make([]byte, 1024)
		dst, err := os.Create("uploaded-file.txt")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
		for {
			n, err := r.Body.Read(buf)
			if err != nil && err != io.EOF {
				return err
			}
			dst.WriteAt(buf[:n], count)
			ntest = (int64)(n)
			count += ntest

			if err == io.EOF {
				break
			}
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return errors.New("Only POST method is supported")
	}

	return nil
}

// UploadForm is a method used to reply user hello in json format
func (r *RestFulUpload) UploadForm(b *rf.Context) {
	// USAGE: curl -X POST http://127.0.0.1:8083/uploadform -H 'content-type: multipart/form-data' -F uploadfile=@input.txt
	req := b.ReadRequest()
	res := b.ReadResponseWriter()
	err := uploadForm(res, req)
	if err == nil {
		b.Write([]byte("Upload multi-form OK!"))
	}
}

func uploadForm(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	//POST takes the uploaded file(s) and saves it to disk.
	case http.MethodPost:
		//parse the multipart form in the request
		err := r.ParseMultipartForm(1024)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}

		//get a ref to the parsed multipart form
		m := r.MultipartForm

		//get the *fileheaders
		files := m.File["uploadfile"]
		for i := range files {
			//for each fileheader, get a handle to the actual file
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return err
			}
			//create destination file making sure the path is writeable.
			dst, err := os.Create("uploaded-form.txt")

			defer dst.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return err
			}
			//copy the uploaded file to the destination file
			if _, err := io.Copy(dst, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return err
			}
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return errors.New("Only POST method is supported")
	}

	return nil
}

//URLPatterns helps to respond for corresponding API calls
func (r *RestFulUpload) URLPatterns() []rf.Route {
	return []rf.Route{
		{http.MethodPost, "/uploadfile", "UploadFile"},
		{http.MethodPost, "/uploadform", "UploadForm"},
	}
}
