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

func (ts *TodoStorage) CreateListRecord(list *TodoListPayload, email string) (int, int, error) {
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
todo_status VARCHAR(10) DEFAULT 'inprogress',
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

func (ts *TodoStorage) GetAllLists(email string) ([]TodoList, error) {
	query := "SELECT id FROM users WHERE email=$1"
	var id int

	records, err := ts.db.Query(query, email)
	if err != nil {
		return []TodoList{}, err
	}

	for records.Next() {
		if err = records.Scan(&id); err != nil {
			return []TodoList{}, err
		}
	}

	query = "SELECT * FROM lists WHERE user_id=$1;"
	var lists []TodoList

	records, err = ts.db.Query(query, id)
	if err != nil {
		return []TodoList{}, err
	}

	for i := 0; records.Next(); i++ {
		var list TodoList
		err = records.Scan(&list.ListID, &list.UserID, &list.Title, &list.CreatedAt)
		if err != nil {
			return []TodoList{}, err
		}
		lists = append(lists, list)
	}

	return lists, nil
}

func (ts *TodoStorage) CreateItemRecord(signature string, payload *TodoItemPayload) (TodoItem, error) {
	query := fmt.Sprintf("INSERT INTO %v(todo_title, todo_description) VALUES($1, $2)", signature)
	_, err := ts.db.Exec(query, payload.Title, payload.Description)
	if err != nil {
		return TodoItem{}, err
	}

	var newItem TodoItem
	query = fmt.Sprintf("SELECT * FROM %v ORDER BY todo_id DESC LIMIT 1;", signature)
	record, err := ts.db.Query(query)
	if err != nil {
		return TodoItem{}, err
	}

	if record.Next() {
		err := record.Scan(&newItem.ID, &newItem.Title, &newItem.Description, &newItem.Status, &newItem.CreatedAt, &newItem.UpdatedAt)
		if err != nil {
			return TodoItem{}, err
		}
	}

	return newItem, nil
}

func (ts *TodoStorage) UpdateItemRecord(signature string, todo_id string, item TodoItemPayload) (TodoItem, error) {
	query := fmt.Sprintf("UPDATE %v SET todo_title='%v', todo_description='%v', todo_status='%v'", signature, item.Title, item.Description, item.Status)
	_, err := ts.db.Exec(query)
	if err != nil {
		return TodoItem{}, err
	}

	var record TodoItem

	query = fmt.Sprintf("SELECT * FROM %v WHERE todo_id=%v", signature, todo_id)
	records, err := ts.db.Query(query)
	if err != nil {
		return TodoItem{}, err
	}

	if records.Next() {
		err = records.Scan(&record.ID, &record.Title, &record.Description, &record.Status, &record.CreatedAt, &record.UpdatedAt)
		if err != nil {
			return TodoItem{}, err
		}
	}

	return record, nil
}

func (ts *TodoStorage) DeleteItemRecord(signature string, todo_id string) error {
	query := fmt.Sprintf("DELETE FROM %v WHERE todo_id=$1", signature)
	_, err := ts.db.Exec(query, todo_id)
	if err != nil {
		return err
	}

	return nil
}

func (ts *TodoStorage) GetAllItems(signature string) ([]TodoItem, error) {
	query := fmt.Sprintf("SELECT * FROM %v;", signature)
	records, err := ts.db.Query(query)
	if err != nil {
		return []TodoItem{}, err
	}

	var items []TodoItem
	for i := 0; records.Next(); i++ {
		var item TodoItem
		err := records.Scan(&item.ID, &item.Title, &item.Description, &item.Status, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return []TodoItem{}, err
		}
		items = append(items, item)
	}

	return items, nil
}

// func (ts *TodoStorage) func() {}
