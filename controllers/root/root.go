package root

import (
	"elo-rating/tools/render"
	"net/http"

	"github.com/gorilla/mux"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	render.RenderRedirect("/lists", w, r)
}

func SetRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/", HandleRoot).Methods(http.MethodGet)

	return r
}
