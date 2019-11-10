package model

import (
	"errors"
)

// Todo from database
type Todo struct {
	ID     string `json:"id"`
	ListID string `json:"list_id"`
	Todo   string `json:"todo"`
	Done   bool   `json:"done"`
}

// TodoList kind of a has many struct
type TodoList struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Todos []*Todo `json:"todos"`
}

// TodoListByID fetches a list and its todos
func TodoListByID(id string) (*TodoList, error) {
	list := &TodoList{}

	sql := "SELECT id, name FROM lists WHERE id = $1"
	err := db.QueryRow(sql, id).Scan(&list.ID, &list.Name)
	if err != nil {
		return nil, err
	}

	list.Todos, err = TodosByListID(id)

	return list, nil
}

func (tl *TodoList) validate() error {
	if tl.Name == "" {
		return errors.New("name empty, must be set")
	}
	if len(tl.Todos) == 0 {
		return errors.New("no todos, must have at least one")
	}

	return nil
}

// Save saves a new list to the database
func (tl *TodoList) Save() error {
	if err := tl.validate(); err != nil {
		return err
	}
	sql := `INSERT INTO lists(name) VALUES($1) RETURNING id;`
	err := db.QueryRow(sql, tl.Name).Scan(&tl.ID)
	if err != nil {
		return err
	}

	for _, t := range tl.Todos {
		t.ListID = tl.ID
		err := t.Save()
		if err != nil {
			return err
		}
	}

	return nil
}

// TodoByID fetches a single todo given a listID and todoID
func TodoByID(listID, todoID string) (*Todo, error) {
	t := &Todo{}
	sql := `SELECT id, list_id, todo, done FROM todos WHERE list_id = $1 AND id = $2 ORDER BY created_at ASC;`
	err := db.QueryRow(sql, listID, todoID).Scan(&t.ID, &t.ListID, &t.Todo, &t.Done)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// TodosByListID fetches the todos for a given list ID
func TodosByListID(id string) ([]*Todo, error) {
	sql := `SELECT id, list_id, todo, done FROM todos WHERE list_id = $1 ORDER BY created_at ASC;`
	rows, err := db.Query(sql, id)
	if err != nil {
		return nil, err
	}

	var todos []*Todo
	for rows.Next() {
		todo := &Todo{}
		err := rows.Scan(&todo.ID, &todo.ListID, &todo.Todo, &todo.Done)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func (t *Todo) validate() error {
	if t.ListID == "" {
		return errors.New("ListID is empty, must be set")
	}
	if t.Todo == "" {
		return errors.New("Todo is empty, must be set")
	}
	return nil
}

// Save saves an individual todo to the database
func (t *Todo) Save() error {
	if err := t.validate(); err != nil {
		return err
	}

	sql := `INSERT INTO todos(list_id, todo) VALUES($1, $2) RETURNING id`
	err := db.QueryRow(sql, t.ListID, t.Todo).Scan(&t.ID)
	if err != nil {
		return err
	}

	return nil
}

// ToggleDone toggles done for the todo and saves to db
func (t *Todo) ToggleDone() error {
	sql := `UPDATE todos SET done = $1 WHERE id = $2 AND list_id = $3 RETURNING done`
	err := db.QueryRow(sql, !t.Done, t.ID, t.ListID).Scan(&t.Done)
	if err != nil {
		return err
	}

	return nil
}
