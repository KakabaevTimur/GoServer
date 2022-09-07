package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const port = "8080"

func main() {
	db := new(DB)

	if len(os.Args) > 1 && os.Args[1] == "-d" {
		db.init(true)
	} else {
		db.init(false)
	}

	defer db.deinit()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" { // curl -d example.com -X POST http://localhost:8080
			b, _ := ioutil.ReadAll(r.Body)
			url := string(b[:])
			io.WriteString(w, "http://localhost:"+port+"/"+db.put(url))
		}
		if r.Method == "GET" {
			req := r.URL.Path[1:]
			io.WriteString(w, db.get(req))
		}
	})

	err := http.ListenAndServe(":"+port, nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
