package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/csrf"
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
	p := csrf.Protect([]byte(os.Getenv("CSRF_TOKEN")))
	err := http.ListenAndServe(":"+port, p(r))
	if err != nil {
		panic(err)
	}
}
