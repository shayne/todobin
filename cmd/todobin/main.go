package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"gitub.com/imartingraham/todobin/internal/route"
)

func main() {
	go route.HandleMessages()
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./web/public"))
	r.PathPrefix("/scripts/").Handler(fs)
	r.PathPrefix("/styles/").Handler(fs)
	r.HandleFunc("/todo/{listId}", route.HandleTodos)
	r.HandleFunc("/todo/{listId}/done/{todoId}", route.HandleTodoDone)
	r.HandleFunc("/ws", route.HandleWs)
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
