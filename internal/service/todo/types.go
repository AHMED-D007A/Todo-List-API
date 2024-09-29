package todo

type TodoListPayload struct {
	Title string `json:"title"`
}

type TodoList struct {
	ListID    int    `json:"list_id"`
	UserID    int    `json:"user_id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
}

type TodoItemPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type TodoItem struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_aT"`
	UpdatedAt   string `json:"updated_at"`
}
