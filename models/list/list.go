package list

import (
	db "elo-rating/database"
	"elo-rating/models/item"
)

type List struct {
	ID    int64
	Title string
}

func All() (*[]List, error) {
	rows, err := db.Lists.All()
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

func Find(id *int64) (*List, error) {
	row := db.Lists.Find(id)

	var list List
	err := row.Scan(&list.ID, &list.Title)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

type ListParams struct {
	Title string
}

func Create(listParams *ListParams) (*int64, error) {
	return db.Lists.Add(&listParams.Title)
}

func (list *List) Delete() error {
	return db.Lists.Delete(&list.ID)
}

func (list *List) Items() (*[]item.Item, error) {
	rows, err := db.Lists.Items(&list.ID)
	if err != nil {
		return nil, err
	}

	return item.DBRowsIntoCollection(rows)
}
