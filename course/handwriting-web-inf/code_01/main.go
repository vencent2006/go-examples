package main

import (
	"go-examples/course/handwriting-web-inf/code_01/framework"
	"log"
	"net/http"
)

func main() {
	server := &http.Server{
		Handler: framework.NewCore(),
		Addr:    "localhost:8080",
	}
	log.Fatal(server.ListenAndServe())
}
