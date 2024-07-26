package routes

import (
	"log"
	"net/http"
	"strconv"

	"github.com/sfluor/test-htmx/views/components"
	"github.com/sfluor/test-htmx/views/model"
)


func RegisterTodos(db model.TodoDB) {
	http.HandleFunc("GET /static/tailwind_out.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/tailwind_out.css")
	})

	http.HandleFunc("GET /todos", func(w http.ResponseWriter, r *http.Request) {
		component := components.Todos(db.ListTodos())
		component.Render(r.Context(), w)
	})

	http.HandleFunc("POST /todos", func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("title")

		if title == "" {
			http.Error(w, "Invalid empty title for todo", 400)
		}

        todo := db.AddTodo(title)
		log.Printf("Created todo; %+v", r.URL.Path)
		components.Todo(todo).Render(r.Context(), w)
	})

	http.HandleFunc("DELETE /todos/{id}", func(w http.ResponseWriter, r *http.Request) {
        rawID := r.PathValue("id")
        ID, err := strconv.Atoi(rawID)
        if err != nil {
            http.Error(w, "Invalid ID received: " + rawID, 400)
        }


        removed, found := db.RemoveTodo(ID)
        if found {
            log.Printf("Deleted todo: %+v", removed)
        } else {
            http.Error(w, "No Todo found with ID: " + rawID, 404)
        }
	})
}
