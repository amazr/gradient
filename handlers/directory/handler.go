package directory

import (
	"example/hello/components"
	"example/hello/services/data"
	"example/hello/services/data/connectors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func GetDirectory(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    dataService := data.NewSqlite()
    view(w, r,dataService.List(1))
}

func GetContent(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    dataService := data.NewSqlite()
    i, err := strconv.Atoi(p.ByName("id"))
    if err != nil {
        panic(err)
    }
    fmt.Fprint(w, string(dataService.Read(i)))
}

func DeleteContent(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    dataService := data.NewSqlite()
    i, err := strconv.Atoi(p.ByName("id"))
    if err != nil {
        panic(err)
    }
    dataService.Delete(i)
    GetDirectory(w, r, p)
}

func Upload(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    dataService := data.NewSqlite()
    f, handler, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Error retrieving file", http.StatusBadRequest)
        return
    }
    defer f.Close()

    dataService.Write(1, handler.Filename, handler.Size, f)
    GetDirectory(w, r, p)
}

func view(w http.ResponseWriter, r *http.Request, locators []connectors.FileLocator) {
    components.Directory(locators).Render(r.Context(), w)
}
