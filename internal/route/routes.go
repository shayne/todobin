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
	"github.com/microcosm-cc/bluemonday"
	"gitub.com/imartingraham/todobin/internal/model"
)

// HandleIndex is the route for "/"
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("web/template/index.html"))
	var tplvars struct {
		Name string
		Todo string
	}
	switch r.Method {
	case "GET":

	case "POST":
		// Disallow all html tags
		p := bluemonday.StrictPolicy()

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		var todos []string
		tplvars.Todo = r.FormValue("todolist")
		listName := p.Sanitize(r.FormValue("name"))

		rawTodos := strings.Split(r.FormValue("todolist"), "\n")

		for _, t := range rawTodos {
			if strings.HasPrefix(t, "-") {
				t = strings.Replace(t, "-", "", 1)
			}
			// When I sanitize the string before splitting
			// it doesn't work, so for now I'm just sanitizing
			// each line
			t = p.Sanitize(strings.TrimSpace(t))
			todos = append(todos, t)
		}

		insertedTodos, err := model.CreateTodos(listName, todos)
		if err != nil {
			panic(err)
		}

		u := "/todo/" + insertedTodos[0].ListID
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
	tpl := template.Must(template.ParseFiles("web/template/todo.html"))
	vars := mux.Vars(r)
	listID, ok := vars["listId"]
	if !ok {
		log.Println("[error] listId not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	list, err := model.GetTodosListByID(listID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = tpl.ExecuteTemplate(w, "todo.html", struct {
		TodoList *model.TodoList
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
		w.WriteHeader(http.StatusForbidden)
		log.Fatalf("[error] Route only handles POST requests")
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("[error] failed to read request body: %v\n", err)
	}

	defer r.Body.Close()

	todo := &model.Todo{}
	err = json.Unmarshal(b, &todo)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("[error] failed to unmarshal json: %v\n", err)
	}

	vars := mux.Vars(r)
	listID := vars["listId"]
	todoID := vars["todoId"]

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
