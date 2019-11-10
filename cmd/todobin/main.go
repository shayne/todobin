package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gitub.com/imartingraham/todobin/internal/route"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/todo/{listId}", route.HandleTodos)
	r.HandleFunc("/todo/{listId}/done/{todoId}", route.HandleTodoDone)
	r.HandleFunc("/", route.HandleIndex)

	http.Handle("/", r)

	fmt.Println("Ready and listening on :3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
