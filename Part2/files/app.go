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

func readiness(w http.ResponseWriter, req *http.Request) {
	resp, err := http.Get("http://readiness.default.svc.cluster.local:8000")
	if err != nil || resp.StatusCode != 200 {
		w.WriteHeader(502)
		fmt.Fprint(w, "Readiness Down")
	} else {
		fmt.Fprintf(w, "Readiness OK\n")
	}
	fmt.Printf("%s %s %s\n", req.Host, req.Method, req.URL)
	log.Printf("%s %s %s\n", req.Host, req.Method, req.URL)
}

func liveness(w http.ResponseWriter, req *http.Request) {
	resp, err := http.Get("http://liveness.default.svc.cluster.local:8000")
	if err != nil || resp.StatusCode != 200 {
		w.WriteHeader(502)
		fmt.Fprint(w, "Liveness Down")
	} else {
		fmt.Fprintf(w, "Liveness OK\n")
	}
	fmt.Printf("%s %s %s\n", req.Host, req.Method, req.URL)
	log.Printf("%s %s %s\n", req.Host, req.Method, req.URL)
}

func main() {
	fileName := "/tmp/go-server.log"
	logFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
	defer logFile.Close()
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/readiness", readiness)
	http.HandleFunc("/liveness", liveness)
	fmt.Println("Listening on 8090")
	http.ListenAndServe(":8090", nil)
}
