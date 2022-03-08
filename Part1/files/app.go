package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
	fmt.Printf("%s %s %s\n", req.Host, req.Method, req.URL)
	log.Printf("%s %s %s\n", req.Host, req.Method, req.URL)
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
	fmt.Printf("%s %s %s\n", req.Host, req.Method, req.URL)
	log.Printf("%s %s %s\n", req.Host, req.Method, req.URL)
}

func main() {
	fileName := "/var/log/go-server.log"
	logFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
	defer logFile.Close()
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	fmt.Println("Listening on 8090")
	http.ListenAndServe(":8090", nil)
}
