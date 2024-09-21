package item

import (
	"database/sql"
	db "elo-rating/database"
)

var defaultRating int = 1000

type Item struct {
	ID     int64
	Title  string
	Rating int
	ListID int64
}

type ItemParams struct {
	Title  string
	ListID int64
}

func DBRowsIntoCollection(rows *sql.Rows) (*[]Item, error) {
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Title, &item.Rating, &item.ListID)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return &items, nil
}

func DBRowIntoItem(row *sql.Row) (*Item, error) {
	var item Item
	err := row.Scan(&item.ID, &item.Title, &item.Rating, &item.ListID)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func Find(id *int64) (*Item, error) {
	row := db.Items.Find(id)
	return DBRowIntoItem(row)
}

func FindForList(id *int64, listID *int64) (*Item, error) {
	row := db.Items.FindForList(id, listID)
	return DBRowIntoItem(row)
}

func Create(itemParams *ItemParams) (*int64, error) {
	id, err := db.Items.Add(&itemParams.Title, &defaultRating, &itemParams.ListID)
	return id, err
}

func (item *Item) Delete() error {
	err := db.Items.Delete(&item.ID)
	return err
}
