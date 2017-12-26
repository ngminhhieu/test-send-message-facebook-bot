package main

import (
	"bkgo/beekit/transport/transhttp"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type test struct {
	Data int64 `json:"url_image"`
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)

	f, handle, err := r.FormFile("file")

	if err != nil {
		transhttp.RespondError(w, http.StatusNotFound, "No such file")
	}
	defer f.Close()

	file, err := os.OpenFile("/home/minhhieu/gov2/src/bkgo/teams/api/happy-messenger-api/"+handle.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("err opening file: %s", err)
	}

	defer file.Close()
	io.Copy(file, f)

	fileInfo, _ := file.Stat()

	size := fileInfo.Size()

	data := &test{Data: size}
	transhttp.RespondJSON(w, http.StatusOK, data)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/upload", ServeHTTP).Methods("POST")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}
