package schemas

import (
	"errors"
	"log"
	"net/http"

	rf "github.com/ServiceComb/go-chassis/server/restful"
	"io"
	"os"
)

//RestFulHello is a struct used for implementation of restfull hello program
type RestFulHello struct {
}

//Sayhello is a method used to reply user with hello
func (r *RestFulHello) Sayhello(b *rf.Context) {
	id := b.ReadPathParameter("userid")
	log.Printf("get user id: " + id)
	b.Write([]byte("get user id: " + id))
}

//Sayhi is a method used to reply user with hello world text
func (r *RestFulHello) Sayhi(b *rf.Context) {
	result := struct {
		Name string
	}{}
	err := b.ReadEntity(&result)
	if err != nil {
		b.Write([]byte(err.Error() + ":hello world"))
		return
	}
	b.Write([]byte(result.Name + ":hello world"))
	return
}

// SayJSON is a method used to reply user hello in json format
func (r *RestFulHello) SayJSON(b *rf.Context) {
	reslut := struct {
		Name string
	}{}
	err := b.ReadEntity(&reslut)
	if err != nil {
		b.WriteHeaderAndJSON(http.StatusInternalServerError, reslut, "application/json")
		return
	}
	reslut.Name = "hello " + reslut.Name
	b.WriteJSON(reslut, "application/json")
	return
}

//URLPatterns helps to respond for corresponding API calls
func (r *RestFulHello) URLPatterns() []rf.Route {
	return []rf.Route{
		{http.MethodGet, "/sayhello/{userid}", "Sayhello"},
		{http.MethodPost, "/sayhi", "Sayhi"},
		{http.MethodPost, "/sayjson", "SayJSON"},
	}
}

//RestFulMessage is a struct used to implement restful message
type RestFulMessage struct {
}

//Saymessage is used to reply user with his name
func (r *RestFulMessage) Saymessage(b *rf.Context) {
	id := b.ReadPathParameter("name")

	b.Write([]byte("get name: " + id))
}

//Sayhi is a method used to reply request user with hello world text
func (r *RestFulMessage) Sayhi(b *rf.Context) {
	reslut := struct {
		Name string
	}{}
	err := b.ReadEntity(&reslut)
	if err != nil {
		b.Write([]byte(err.Error() + ":hello world"))
		return
	}
	b.Write([]byte(reslut.Name + ":hello world"))
	return
}

//Sayerror is a method used to reply request user with error
func (r *RestFulMessage) Sayerror(b *rf.Context) {
	b.WriteError(http.StatusInternalServerError, errors.New("test hystric"))
	return
}

//URLPatterns helps to respond for corresponding API calls
func (r *RestFulMessage) URLPatterns() []rf.Route {
	return []rf.Route{
		{http.MethodGet, "/saymessage/{name}", "Saymessage"},
		{http.MethodPost, "/sayhimessage", "Sayhi"},
		{http.MethodGet, "/sayerror", "Sayerror"},
	}
}

//RestFulUpload is a struct used to implement restful upload of Multiform and raw file data
type RestFulUpload struct {
}

//UploadFile is used to upload a binary file
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
		dst, err := os.Create("test.txt")
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

//UploadMultiForm is used to upload a binary file
func (r *RestFulUpload) UploadMultiForm(b *rf.Context) {
	// USAGE: curl -X POST http://127.0.0.1:8083/uploadmultiform -H 'content-type: multipart/form-data' -F uploadfile=@input.txt
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
			dst, err := os.Create("output_" + files[i].Filename)

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
		{http.MethodPost, "/uploadmultiform", "UploadMultiForm"},
	}
}
