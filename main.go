package main

import (
	"log"
	"net/http"

	"github.com/sfluor/test-htmx/views/model"
	"github.com/sfluor/test-htmx/views/model/buldan"
	"github.com/sfluor/test-htmx/views/routes"
)

func main() {

	http.HandleFunc("GET /static/tailwind_out.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/tailwind_out.css")
	})

    routes.RegisterTodos(model.NewTodoDB())
    routes.RegisterBuldan(buldan.NewEngine())

	log.Print("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
