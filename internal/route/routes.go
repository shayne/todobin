package route

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"gitub.com/imartingraham/todobin/internal/model"
)

// HandleIndex is the route for "/"
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("web/template/index.html"))
	var tplvars struct {
		CSRFToken string
		Name      string
		Todo      string
	}

	switch r.Method {
	case "GET":
		tplvars.CSRFToken = csrf.Token(r)

	case "POST":
		// Disallow all html tags
		p := bluemonday.StrictPolicy()

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		todoList := &model.TodoList{
			Name: p.Sanitize(r.FormValue("name")),
		}

		tplvars.Todo = r.FormValue("todolist")

		rawTodos := strings.Split(r.FormValue("todolist"), "\n")

		for _, t := range rawTodos {
			if strings.HasPrefix(t, "-") {
				t = strings.Replace(t, "-", "", 1)
			}
			// When I sanitize the string before splitting
			// it doesn't work, so for now I'm just sanitizing
			// each line
			todoList.Todos = append(todoList.Todos, &model.Todo{
				Todo: p.Sanitize(strings.TrimSpace(t)),
			})
		}

		err := todoList.Save()
		if err != nil {
			panic(err)
		}

		u := "/todo/" + todoList.ID
		http.Redirect(w, r, u, 301)
		return
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
	err := tpl.ExecuteTemplate(w, "index.html", tplvars)
	if err != nil {
		log.Fatalf("[error] failed to execute index.html: %v\n", err)
	}
}

// HandleTodos is the route for "/todo/[uuid]"
func HandleTodos(w http.ResponseWriter, r *http.Request) {

	var tplvars struct {
		CSRFToken string
		TodoList  *model.TodoList
	}

	tplvars.CSRFToken = csrf.Token(r)
	tpl := template.Must(template.ParseFiles("web/template/todo.html"))
	vars := mux.Vars(r)
	listID, ok := vars["listId"]
	if !ok {
		log.Println("[error] listId not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	list, err := model.TodoListByID(listID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	tplvars.TodoList = list
	err = tpl.ExecuteTemplate(w, "todo.html", tplvars)
	if err != nil {
		log.Fatalf("[error] failed to execute todo.html: %v\n", err)
	}
}

// HandleTodoDone marks todo as done or undone
func HandleTodoDone(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusForbidden)
		log.Fatalf("[error] Route only handles PUT requests")
		return
	}

	vars := mux.Vars(r)
	listID := vars["listId"]
	todoID := vars["todoId"]

	t, err := model.TodoByID(listID, todoID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("[error] failed to fetch todo: %v\n", err)
	}

	err = t.ToggleDone()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("[error] failed to toggle done for todo: %v\n", err)
	}

	jsonData, err := json.Marshal(t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("[error] could not json encode todo: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
