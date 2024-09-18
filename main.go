package main

// TODO: Add table names as constants

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Item struct {
	ID     int
	Title  string
	Rating int
	ListID int
}

type List struct {
	ID    int
	Title string
	Items *[]Item
}

type ListsPageData struct {
	Lists *[]List
}

type ItemsPageData struct {
	ListID string
	List   *List
	Items  *[]Item
}

type Database struct {
	db *sql.DB
}

func (database *Database) createListsTable() error {
	_, err := database.db.Exec(`
		CREATE TABLE IF NOT EXISTS lists (
		id INTEGER PRIMARY KEY,
		title TEXT
		)`)
	return err
}

func (database *Database) createItemsTable() error {
	_, err := database.db.Exec(`
		CREATE TABLE IF NOT EXISTS items (
		id INTEGER PRIMARY KEY,
		title TEXT,
		rating INTEGER,
		list_id INTEGER,
		FOREIGN KEY (list_id) REFERENCES lists(id)
		)`)
	return err
}

func (database *Database) getLists() (*[]List, error) {
	rows, err := database.db.Query("SELECT id, title FROM lists")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lists []List
	for rows.Next() {
		var list List
		err = rows.Scan(&list.ID, &list.Title)
		if err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}

	return &lists, nil
}

func (database *Database) getList(id string) (*List, error) {
	row := database.db.QueryRow("SELECT title FROM lists WHERE id = ?", id)
	var title string
	err := row.Scan(&title)

	if err != nil {
		return nil, err
	}

	list := List{
		Title: title,
	}

	return &list, nil
}

func (database *Database) deleteList(id string) error {
	_, err := database.db.Exec("DELETE FROM lists WHERE id = ?", id)
	return err
}

type ListParams struct {
	title string
}

func (database *Database) addList(list ListParams) error {
	_, err := database.db.Exec("INSERT INTO lists (title) VALUES (?)", list.title)
	return err
}

func (database *Database) getItemsForList(listID string) (*[]Item, error) {
	rows, err := database.db.Query("SELECT id, title, rating, list_id FROM items WHERE list_id = ?", listID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err = rows.Scan(&item.ID, &item.Title, &item.Rating, &item.ListID)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return &items, nil
}

func (database *Database) deleteItemForList(id string, listID string) error {
	_, err := database.db.Exec("DELETE FROM items WHERE id = ? AND list_id = ?", id, listID)
	return err
}

type ItemParams struct {
	title  string
	listID string
}

func (database *Database) addItem(item ItemParams) error {
	defaultRating := 1000
	_, err := database.db.Exec("INSERT INTO items (title, rating, list_id) VALUES (?, ?, ?)", item.title, defaultRating, item.listID)
	return err
}

func (database *Database) migrate() {
	err := database.createListsTable()
	if err != nil {
		log.Fatal(err)
	}

	err = database.createItemsTable()
	if err != nil {
		log.Fatal(err)
	}
}

func InitDatabase(filePath string) *Database {
	// Initialize the SQLite database
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		log.Fatal(err)
	}

	return &Database{db: db}
}

func renderTemplate(templatePath string, w http.ResponseWriter, args any) error {
	tmpl := template.Must(template.ParseFiles(templatePath))
	return tmpl.Execute(w, args)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/lists", http.StatusSeeOther)
}

func handleLists(database *Database) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			lists, err := database.getLists()
			if err != nil {
				log.Fatal(err)
			}

			listPageData := ListsPageData{
				Lists: lists,
			}

			err = renderTemplate("templates/lists/index.html", w, listPageData)
			if err != nil {
				log.Fatal(err)
			}
		case http.MethodPost:
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}

			title := r.FormValue("title")

			err = database.addList(ListParams{title: title})
			if err != nil {
				log.Fatal(err)
			}

			http.Redirect(w, r, "/lists", http.StatusSeeOther)
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
	return http.HandlerFunc(fn)
}

func handleList(database *Database) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id := vars["id"]
		if id == "" {
			http.Error(w, "Missing ID", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodPost:
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}

			method := r.FormValue("_method")
			// Can't send Delete from plain HTML
			if method == "delete" {

				err := database.deleteList(id)
				if err != nil {
					log.Fatal(err)
				}

				// Redirect back to the homelist
				http.Redirect(w, r, "/lists", http.StatusSeeOther)
			} else {
				http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			}
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
	return http.HandlerFunc(fn)
}

func handleItems(database *Database) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		listID := vars["listID"]
		if listID == "" {
			http.Error(w, "Missing list ID", http.StatusBadRequest)
			return
		}

		list, err := database.getList(listID)
		if err != nil {
			http.Error(w, "List not found. List ID: "+listID, http.StatusNotFound)
			return
		}

		switch r.Method {
		case http.MethodGet:
			items, err := database.getItemsForList(listID)
			if err != nil {
				log.Fatal(err)
			}

			itemsPageData := ItemsPageData{
				ListID: listID,
				List:   list,
				Items:  items,
			}

			err = renderTemplate("templates/lists/items/index.html", w, itemsPageData)
			if err != nil {
				log.Fatal(err)
			}
		case http.MethodPost:
			title := r.FormValue("title")

			err := database.addItem(ItemParams{title: title, listID: listID})
			if err != nil {
				log.Fatal(err)
			}

			http.Redirect(w, r, "/lists/"+listID+"/items", http.StatusSeeOther)
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}

	return http.HandlerFunc(fn)
}

func handleItem(database *Database) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		listID := vars["listID"]
		if listID == "" {
			http.Error(w, "Missing List ID", http.StatusBadRequest)
			return
		}

		id := vars["id"]
		if id == "" {
			http.Error(w, "Missing ID", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodPost:
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}

			method := r.FormValue("_method")
			// Can't send Delete from plain HTML
			if method == "delete" {
				// Delete the list from the database
				err := database.deleteItemForList(id, listID)
				if err != nil {
					log.Fatal(err)
				}

				// Redirect back to the homelist
				http.Redirect(w, r, "/lists/"+listID+"/items", http.StatusSeeOther)
			} else {

				http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			}
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
	return http.HandlerFunc(fn)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	database := InitDatabase("./app.db")
	database.migrate()
	defer database.db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/", handleRoot).Methods(http.MethodGet)
	r.HandleFunc("/lists", handleLists(database)).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/lists/{id:[0-9]+}", handleList(database)).Methods(http.MethodPost)
	r.HandleFunc("/lists/{listID:[0-9]+}/items", handleItems(database)).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/lists/{listID:[0-9]+}/items/{id:[0-9]+}", handleItem(database)).Methods(http.MethodPost)

	//
	// listsRouter := r.PathPrefix("/lists").Subrouter()
	// listsRouter.HandleFunc("/", handleLists(database)).Methods(http.MethodGet, http.MethodPost)
	// listsRouter.HandleFunc("/{id:[0-9]+}", handleList(database)).Methods(http.MethodGet, http.MethodDelete)
	//
	// itemsRouter := listsRouter.PathPrefix("/{listID:[0-9]+}/items").Subrouter()
	// itemsRouter.HandleFunc("/", handleItems(database)).Methods(http.MethodGet, http.MethodPost)
	// itemsRouter.HandleFunc("/{id:[0-9]+}", handleItem(database)).Methods(http.MethodGet, http.MethodDelete)
	//
	http.Handle("/", r)

	// Start the server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
