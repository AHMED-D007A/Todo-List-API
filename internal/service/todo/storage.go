package todo

import "database/sql"

type TodoStorage struct {
	db *sql.DB
}

func NewTodoStorage(db *sql.DB) *TodoStorage {
	return &TodoStorage{
		db: db,
	}
}

func (ts *TodoStorage) CreateListRecord(list *TodoList, email string) error {
	return nil
}

// func (ts *TodoStorage) func() {}
