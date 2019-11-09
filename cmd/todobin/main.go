package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitub.com/imartingraham/todobin/internal/model"
	"gitub.com/imartingraham/todobin/internal/route"

	_ "github.com/lib/pq"
)

func main() {
	model.InitDB()
	r := mux.NewRouter()
	r.HandleFunc("/todo/{listId}", route.HandleTodos)
	r.HandleFunc("/todo/{listId}/done/{todoId}", route.HandleTodoDone)
	r.HandleFunc("/", route.HandleIndex)

	http.Handle("/", r)
	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		panic(err)
	}
}
