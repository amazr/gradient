package main

import (
	"example/hello/handlers"
	"log"
	"net/http"
) 

func main() {
    //b := lang.Ser(lang.New("632a1792-1ce3-4de6-bf0f-aea5defdcfb5", lang.NewKeepFilter(lang.NewIsEqual("Type 1", "Grass"))))
    //r := lang.Des(b)
    //e := exec.NewSqlite3Executor()
    //result := e.Execute(r)
    //fmt.Println(result)

    log.Fatal(http.ListenAndServe("127.0.0.1:8000", handlers.BuildRouter()));
}
