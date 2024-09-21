package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

var Connection *Database

func newConnection(db *sql.DB) *Database {
	return &Database{
		db: db,
	}
}

func InitDatabase(filePath string) *Database {
	// Initialize the SQLite database
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		log.Fatal(err)
	}

	Connection = newConnection(db)
	return Connection
}

func (db *Database) Close() {
	db.db.Close()
}

func (*Database) Migrate() {
	err := Lists.createTable()
	if err != nil {
		log.Fatal(err)
	}

	err = Items.createTable()
	if err != nil {
		log.Fatal(err)
	}
}

// lists table
type ListsModel struct{}

var Lists *ListsModel

func (*ListsModel) createTable() error {
	_, err := Connection.db.Exec(`
		CREATE TABLE IF NOT EXISTS lists (
		id INTEGER PRIMARY KEY,
		title TEXT
		)`)
	return err
}

func (*ListsModel) All() (*sql.Rows, error) {
	return Connection.db.Query("SELECT id, title FROM lists")
}

func (*ListsModel) Find(id *int64) *sql.Row {
	return Connection.db.QueryRow("SELECT id, title FROM lists WHERE id = ?", *id)
}

func (*ListsModel) Add(title *string) (*int64, error) {
	res, err := Connection.db.Exec("INSERT INTO lists (title) VALUES (?)", title)
	if err != nil {
		return nil, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &lastInsertID, err
}

func (*ListsModel) Delete(id *int64) error {
	_, err := Connection.db.Exec("DELETE FROM lists WHERE id = ?", id)
	return err
}

func (*ListsModel) Items(listID *int64) (*sql.Rows, error) {
	return Connection.db.Query("SELECT id, title, rating, list_id FROM items WHERE list_id = ?", *listID)
}

// items table

type ItemsModel struct{}

var Items *ItemsModel

func (*ItemsModel) createTable() error {
	_, err := Connection.db.Exec(`
		CREATE TABLE IF NOT EXISTS items (
		id INTEGER PRIMARY KEY,
		title TEXT,
		rating INTEGER,
		list_id INTEGER,
		FOREIGN KEY (list_id) REFERENCES lists(id)
		)`)
	return err
}

func (*ItemsModel) Find(id *int64) *sql.Row {
	return Connection.db.QueryRow("SELECT id, title, rating, list_id FROM items WHERE id = ?", *id)
}

func (*ItemsModel) FindForList(id *int64, listID *int64) *sql.Row {
	return Connection.db.QueryRow("SELECT id, title, rating, list_id FROM items WHERE id = ? AND list_id = ?", *id, *listID)
}

func (*ItemsModel) Add(title *string, rating *int, listID *int64) (*int64, error) {
	res, err := Connection.db.Exec("INSERT INTO items (title, rating, list_id) VALUES (?, ?, ?)", *title, *rating, *listID)
	if err != nil {
		return nil, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &lastInsertID, err
}

func (*ItemsModel) Delete(id *int64) error {
	_, err := Connection.db.Exec("DELETE FROM items WHERE id = ?", *id)
	return err
}
