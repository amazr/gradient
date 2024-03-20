package handlers

import (
	"example/hello/components"
	"example/hello/handlers/directory"
	"example/hello/handlers/preview"
	"example/hello/handlers/query"
	"example/hello/services/data"
	"fmt"
	"net/http"

	httprouter "github.com/julienschmidt/httprouter"
)

func BuildRouter() *httprouter.Router {
    router := httprouter.New()

    ds := data.NewSqlite()
    dh := directory.New(ds)
    ph := preview.New(ds)
    qh := query.New(ds)

    router.GET("/", Middleware(dh.GetIndex))
    router.GET("/directory", Middleware(dh.GetFileList))
    router.GET(fmt.Sprintf("/content/:%s", preview.ID_PARAM), Middleware(ph.GetPreview))
    router.DELETE(fmt.Sprintf("/content/:%s", directory.ID_PARAM), Middleware(dh.DeleteFile))
    router.GET(fmt.Sprintf("/board/:%s", query.ID_PARAM), Middleware(qh.GetBoard))
    router.POST(fmt.Sprintf("/filter/:%s", query.ID_PARAM), Middleware(qh.ExecFilter))
    router.POST("/upload", Middleware(dh.Upload))
    defer fmt.Println("Router setup!")

    return router
}

func Middleware(next httprouter.Handle) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        components.Head().Render(r.Context(), w)
        next(w, r, ps)
        components.Scripts().Render(r.Context(), w)
    }
}
