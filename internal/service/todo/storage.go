package todo

import (
	"database/sql"
	"fmt"
)

type TodoStorage struct {
	db *sql.DB
}

func NewTodoStorage(db *sql.DB) *TodoStorage {
	return &TodoStorage{
		db: db,
	}
}

func (ts *TodoStorage) CreateListRecord(list *TodoList, email string) (int, int, error) {
	query := "SELECT id FROM users WHERE email=$1"
	var id int

	records, err := ts.db.Query(query, email)
	if err != nil {
		return 0, 0, err
	}

	for records.Next() {
		if err = records.Scan(&id); err != nil {
			return 0, 0, err
		}
	}

	query = "INSERT INTO lists(user_id, list_title) VALUES($1, $2)"
	_, err = ts.db.Exec(query, id, list.Title)
	if err != nil {
		return 0, 0, err
	}

	query = "SELECT list_id FROM lists WHERE user_id=$1"
	var list_id int

	records, err = ts.db.Query(query, id)
	if err != nil {
		return 0, 0, err
	}

	for records.Next() {
		if err = records.Scan(&list_id); err != nil {
			return 0, 0, err
		}
	}

	query = fmt.Sprintf(`CREATE TABLE list_%v_%v(
todo_id SERIAL PRIMARY KEY,
todo_title VARCHAR(100) NOT NULL,
todo_description TEXT,
created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);`, id, list_id)

	_, err = ts.db.Exec(query)
	if err != nil {
		return 0, 0, err
	}

	query = `CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;`

	_, err = ts.db.Exec(query)
	if err != nil {
		return 0, 0, err
	}

	query = fmt.Sprintf(`CREATE TRIGGER trigger_update_updated_at
BEFORE UPDATE ON list_%v_%v
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();`, id, list_id)

	_, err = ts.db.Exec(query)
	if err != nil {
		return 0, 0, err
	}

	return id, list_id, nil
}

// func (ts *TodoStorage) func() {}
