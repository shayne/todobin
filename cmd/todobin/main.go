package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gitub.com/imartingraham/todobin/internal/route"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/todo/{listId}", route.HandleTodos)
	r.HandleFunc("/todo/{listId}/done/{todoId}", route.HandleTodoDone)
	r.HandleFunc("/", route.HandleIndex)

	http.Handle("/", r)
	port := os.Getenv("PORT")
	fmt.Println("Ready and listening on " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
