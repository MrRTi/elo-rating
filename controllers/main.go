package controllers

import (
	"elo-rating/controllers/items"
	"elo-rating/controllers/lists"
	"elo-rating/controllers/root"
	"net/http"

	"github.com/gorilla/mux"
)

func SetRoutes(r *mux.Router) *mux.Router {
	root.SetRoutes(r)
	listsRouter := lists.SetRoutes(r)
	items.SetRoutes(listsRouter)

	// r.HandleFunc("/", root.HandleRoot).Methods(http.MethodGet)
	// r.HandleFunc("/lists", lists.HandleLists).Methods(http.MethodGet, http.MethodPost)
	// r.HandleFunc("/lists/{id:[0-9]+}", lists.HandleList).Methods(http.MethodPost)
	// r.HandleFunc("/lists/{listID:[0-9]+}/items", items.HandleItems).Methods(http.MethodGet, http.MethodPost)
	// r.HandleFunc("/lists/{listID:[0-9]+}/items/{id:[0-9]+}", items.HandleItem).Methods(http.MethodPost)

	http.Handle("/", r)

	return r
}
