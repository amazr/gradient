package preview

import (
	"example/hello/components/preview"
	"example/hello/services/data"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const ID_PARAM = "id"

type PreviewHandler struct {
    ds data.DataService
}

func New(ds data.DataService) (*PreviewHandler) {
    return &PreviewHandler{
    	ds: ds,
    }
}

func (h *PreviewHandler) GetPreview(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    fid := p.ByName(ID_PARAM)
    cols, rows := h.ds.Read(fid)
    preview.Preview(false, fid, cols, rows).Render(r.Context(), w)
}
