package model

// Todo from database
type Todo struct {
	ID     string `json:"id"`
	ListID string `json:"list_id"`
	Todo   string `json:"todo"`
	Done   bool   `json:"done"`
}

// TodoList kind of a has many struct
type TodoList struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Todos []Todo `json:"todos"`
}

// GetTodosListByID returns slice of Todos based on the listID passed in
func GetTodosListByID(listID string) (TodoList, error) {
	listStatement := `SELECT id, name FROM lists WHERE id = $1`
	list := TodoList{}

	err := db.QueryRow(listStatement, listID).Scan(&list.ID, &list.Name)
	if err != nil {
		return list, err
	}

	sqlStatement := `SELECT id, list_id, todo, done FROM todos WHERE list_id = $1 ORDER BY created_at ASC;`

	rows, err := db.Query(sqlStatement, listID)

	if err != nil {
		return list, err
	}

	for rows.Next() {
		todo := Todo{}
		err := rows.Scan(&todo.ID, &todo.ListID, &todo.Todo, &todo.Done)
		if err != nil {
			return list, err
		}
		list.Todos = append(list.Todos, todo)
	}
	return list, nil
}

// CreateTodos inserts todos into the todos table
func CreateTodos(name string, todos []string) ([]Todo, error) {
	var newTodos []Todo
	listID, err := createList(name)

	if err != nil {
		return nil, err
	}

	sqlStatement := `INSERT INTO todos(list_id, todo) VALUES($1, $2)
		RETURNING id, list_id, todo, done`
	for _, todo := range todos {
		var createdTodo = Todo{}
		err := db.QueryRow(sqlStatement, listID, todo).Scan(&createdTodo.ID, &createdTodo.ListID, &createdTodo.Todo, &createdTodo.Done)
		if err != nil {
			return nil, err
		} else if createdTodo.ID != "" {
			newTodos = append(newTodos, createdTodo)
		}
	}
	return newTodos, nil
}

// createList creates a list and returns the UUID
func createList(name string) (string, error) {
	var id string
	sqlStatement := `INSERT INTO lists(name) VALUES($1) RETURNING id;`
	err := db.QueryRow(sqlStatement, name).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

// MarkTodoAsDone mark todo as done or undone depending on the value passed in
func MarkTodoAsDone(listID string, todoID string, done bool) (Todo, error) {
	todo := Todo{}

	sqlStatement := `UPDATE todos SET done = $1 WHERE id = $2 AND list_id = $3 RETURNING id, list_id, todo, done`
	err := db.QueryRow(sqlStatement, done, todoID, listID).Scan(&todo.ID, &todo.ListID, &todo.Todo, &todo.Done)

	if err != nil {
		return todo, err
	}

	return todo, nil
}
