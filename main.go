package main

import (
	"example/hello/handlers"
	"log"
	"net/http"
) 

func main() {
    log.Fatal(http.ListenAndServe("127.0.0.1:8000", handlers.BuildRouter()));
}
