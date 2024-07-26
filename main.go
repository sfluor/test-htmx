package main

import (
	"log"
	"net/http"

	"github.com/sfluor/test-htmx/views/model"
	"github.com/sfluor/test-htmx/views/routes"
)

func main() {
    todosDB := model.NewTodoDB()
    routes.RegisterTodos(todosDB)

	log.Print("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
