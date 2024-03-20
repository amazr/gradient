package directory

import (
	"example/hello/components/directory"
	"example/hello/services/data"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const ID_PARAM = "id"
const FORM_FILE_PARAM = "file"

type DirectoryHandler struct {
    ds data.DataService
}

func New(ds data.DataService) (*DirectoryHandler) {
    return &DirectoryHandler{
    	ds: ds,
    }
}

func (h *DirectoryHandler) GetIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    directory.Directory(FORM_FILE_PARAM).Render(r.Context(), w)
}

func (h *DirectoryHandler) GetFileList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    directory.DirectoryPostLoad(h.ds.List(1)).Render(r.Context(), w)
}

func (h *DirectoryHandler) DeleteFile(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    h.ds.Delete(p.ByName(ID_PARAM))
    h.GetFileList(w, r, p)
}

func (h *DirectoryHandler) Upload(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    f, handler, err := r.FormFile(FORM_FILE_PARAM)
    defer h.GetFileList(w, r, p)
    if err != nil {
        return
    }
    defer f.Close()

    h.ds.Write(1, handler.Filename, handler.Size, f)
}
