package query

import (
	"example/hello/components/preview"
	"example/hello/components/query"
	"example/hello/services/data"
	"example/hello/services/query/exec"
	"example/hello/services/query/lang"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const ID_PARAM = "id"

type QueryHandler struct {
    ds data.DataService
}

func New(ds data.DataService) (*QueryHandler) {
    return &QueryHandler{
        ds: ds,
    }
}

func (h *QueryHandler) GetBoard(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    fid := p.ByName(ID_PARAM)
    // q := lang.New("632a1792-1ce3-4de6-bf0f-aea5defdcfb5", lang.NewKeepFilter(lang.NewIsEqual("Type 1", "Grass"))))
    //r := lang.Des(b)
    //e := exec.NewSqlite3Executor()
    //result := e.Execute(r)
    //fmt.Println(result)
    cols, rows := h.ds.Read(fid)
    query.Board(fid, cols, rows).Render(r.Context(), w)
}

func (h *QueryHandler) ExecFilter(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    fid := p.ByName(ID_PARAM)
    ftype := r.FormValue("ftype")
    col := r.FormValue("col")
    cond := r.FormValue("cond")
    val := r.FormValue("val")

    var condition lang.Condition
    if cond == "ie" {
        condition = lang.NewIsEqual(col, val)
    } else {
        panic("Unknown cond type")
    }

    var filter lang.Query
    if ftype == "kf" {
        filter = lang.NewKeepFilter(condition)
    } else if ftype == "rf" {
        filter = lang.NewRemoveFilter(condition)
    } else {
        panic("Unknown filter type")
    }

    q := lang.New(fid, filter)
    e := exec.NewSqlite3Executor()
    cols, rows := e.Execute(q)

    preview.Preview(true, fid, cols, rows).Render(r.Context(), w)
}

