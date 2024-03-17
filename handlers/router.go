package handlers

import (
	"example/hello/handlers/directory"
	"example/hello/handlers/index"

	httprouter "github.com/julienschmidt/httprouter"
)

func BuildRouter() *httprouter.Router {
    router := httprouter.New()
    router.GET("/", index.Get)
    router.GET("/directory", directory.GetDirectory)
    router.GET("/content/:id", directory.GetContent)
    router.DELETE("/content/:id", directory.DeleteContent)
    router.POST("/upload", directory.Upload)
    return router
}
