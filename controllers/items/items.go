package items

import (
	"elo-rating/models/item"
	"elo-rating/models/list"
	"elo-rating/tools/render"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ItemsPageData struct {
	ListID int64
	List   *list.List
	Items  *[]item.Item
}

func HandleItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	listIDStr := vars["listID"]
	listID, err := strconv.ParseInt(listIDStr, 10, 64)
	if err != nil || listIDStr == "" {
		render.RenderBadRequest("Missing list ID", w)
		return
	}

	listRecord, err := list.Find(&listID)
	if err != nil {
		errorMessage := fmt.Sprintf("List not found. List ID: %d", listID)
		render.RenderNotFound(errorMessage, w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		items, err := listRecord.Items()
		if err != nil {
			log.Fatal(err)
		}

		itemsPageData := ItemsPageData{
			ListID: listID,
			List:   listRecord,
			Items:  items,
		}

		err = render.RenderTemplate("templates/lists/items/index.html", w, itemsPageData)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodPost:
		title := r.FormValue("title")

		_, err := item.Create(&item.ItemParams{Title: title, ListID: listID})
		if err != nil {
			log.Fatal(err)
		}

		redirectUrl := fmt.Sprintf("/lists/%d/items", listID)
		render.RenderRedirect(redirectUrl, w, r)
	default:
		render.RenderMethodNotAllowed("Invalid request method", w)
	}
}

func HandleItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	listIDStr := vars["listID"]
	listID, err := strconv.ParseInt(listIDStr, 10, 64)
	if err != nil || listIDStr == "" {
		render.RenderBadRequest("Missing list ID", w)
		return
	}

	_, err = list.Find(&listID)
	if err != nil {
		errorMessage := fmt.Sprintf("List not found. List ID: %d", listID)
		render.RenderNotFound(errorMessage, w)
		return
	}

	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || idStr == "" {
		render.RenderBadRequest("Missing ID", w)
		return
	}

	itemRecord, err := item.FindForList(&id, &listID)
	if err != nil {
		errorMessage := fmt.Sprintf("List not found. List ID: %d", listID)
		render.RenderNotFound(errorMessage, w)
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
			err := itemRecord.Delete()
			if err != nil {
				log.Fatal(err)
			}

			redirectUrl := fmt.Sprintf("/lists/%d/items", listID)
			render.RenderRedirect(redirectUrl, w, r)
		} else {

			render.RenderMethodNotAllowed("Invalid request method", w)
		}
	default:
		render.RenderMethodNotAllowed("Invalid request method", w)
	}
}

func SetRoutes(r *mux.Router) *mux.Router {
	itemsRouter := r.PathPrefix("/{listID:[0-9]+}/items").Subrouter()
	itemsRouter.HandleFunc("", HandleItems).Methods(http.MethodGet, http.MethodPost)
	itemsRouter.HandleFunc("/{id:[0-9]+}", HandleItem).Methods(http.MethodPost)

	return itemsRouter
}
