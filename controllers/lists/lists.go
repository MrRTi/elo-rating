package lists

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"elo-rating/models/list"
	"elo-rating/tools/render"
)

type ListsPageData struct {
	Lists *[]list.List
}

func HandleLists(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		lists, err := list.All()
		if err != nil {
			log.Fatal(err)
		}

		listPageData := ListsPageData{
			Lists: lists,
		}

		err = render.RenderTemplate("templates/lists/index.html", w, listPageData)
		if err != nil {
			log.Fatal(err)
		}

	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			render.RenderBadRequest("Failed to parse form", w)
			return
		}

		title := r.FormValue("title")

		_, err = list.Create(&list.ListParams{Title: title})
		if err != nil {
			log.Fatal(err)
		}

		render.RenderRedirect("/lists", w, r)
	default:
		render.RenderMethodNotAllowed("Invalid request method", w)
	}
}

func HandleList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || idStr == "" {

		render.RenderBadRequest("Missing ID", w)
		return
	}

	switch r.Method {
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			render.RenderBadRequest("Failed to parse form", w)
			return
		}

		method := r.FormValue("_method")
		// Can't send Delete from plain HTML
		if method == "delete" {
			listRecord, err := list.Find(&id)
			if err != nil {
				render.RenderNotFound("Not found", w)
			}

			err = listRecord.Delete()
			if err != nil {
				log.Fatal(err)
			}

			render.RenderRedirect("/lists", w, r)
		} else {
			render.RenderMethodNotAllowed("Invalid request method", w)
		}
	default:
		render.RenderMethodNotAllowed("Invalid request method", w)
	}
}

func SetRoutes(r *mux.Router) *mux.Router {
	listsRouter := r.PathPrefix("/lists").Subrouter()
	listsRouter.HandleFunc("", HandleLists).Methods(http.MethodGet, http.MethodPost)
	listsRouter.HandleFunc("/{id:[0-9]+}", HandleList).Methods(http.MethodPost)

	return listsRouter
}
