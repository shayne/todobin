package route

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"gitub.com/imartingraham/todobin/internal/model"
)

// HandleIndex is the route for "/"
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("web/template/index.html"))
	var todo string
	var todos []string

	switch r.Method {
	case "GET":
		todo = ""
		todos = []string{}
	case "POST":

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
		}
		todo = r.FormValue("todolist")
		listName := r.FormValue("name")
		rawTodos := strings.Split(r.FormValue("todolist"), "\n")
		for _, t := range rawTodos {
			if strings.HasPrefix(t, "-") {
				t = strings.Replace(t, "-", "", 1)
			}
			t = strings.Trim(t, " ")
			todos = append(todos, t)
		}

		insertedTodos, err := model.CreateTodos(listName, todos)
		if err != nil {
			panic(err)
		}

		u := fmt.Sprintf("/todo/%s", insertedTodos[0].ListID)
		http.Redirect(w, r, u, 301)
		return
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
	err := tpl.ExecuteTemplate(w, "index.html", struct {
		Todo  string
		Todos []string
	}{
		Todo:  todo,
		Todos: todos,
	})
	if err != nil {
		log.Fatalf("[error] failed to execute index.html: %v\n", err)
	}
}

// HandleTodos is the route for "/todo/[uuid]"
func HandleTodos(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("web/template/todo.html"))
	vars := mux.Vars(r)
	listID := vars["listId"]
	list, err := model.GetTodosListByID(listID)
	if err != nil {
		log.Fatalf("[error] failed to get rows by list id: %v\n", err)
	}
	err = tpl.ExecuteTemplate(w, "todo.html", struct {
		TodoList model.TodoList
	}{
		TodoList: list,
	})
	if err != nil {
		log.Fatalf("[error] failed to execute todo.html: %v\n", err)
	}
}

// HandleTodoDone marks todo as done or undone
func HandleTodoDone(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Fatalf("[error] Route only handles POST requests")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	vars := mux.Vars(r)
	listID := vars["listId"]
	todoID := vars["todoId"]

	todo := model.Todo{}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("[error] failed to read request body: %v\n", err)
	}

	defer r.Body.Close()
	err = json.Unmarshal(b, &todo)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("[error] failed to unmarshal json: %v\n", err)
	}
	todo, err = model.MarkTodoAsDone(listID, todoID, todo.Done)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("[error] failed to update todo: %v\n", err)
	}
	jsonData, err := json.Marshal(todo)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("[error] could not json encode todo: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
