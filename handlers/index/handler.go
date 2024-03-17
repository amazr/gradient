package index

import (
	"example/hello/components"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    components.Index().Render(r.Context(), w)
}
