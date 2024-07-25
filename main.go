package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/sfluor/test-htmx/views/components"
	"github.com/sfluor/test-htmx/views/model"
)

func main() {

	nextID := 4
	todos := []model.Todo{
		{ID: 1, Title: "Clean the dishes"},
		{ID: 2, Title: "Learn HTMX"},
		{ID: 3, Title: "Workout"},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		component := components.Index(todos)
		component.Render(r.Context(), w)
	})

	http.HandleFunc("POST /todos", func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("title")

		if title == "" {
			http.Error(w, "Invalid empty title for todo", 400)
		}

		todo := model.Todo{
			ID:    nextID,
			Title: title,
		}
		todos = append(todos, todo)
		nextID++
		log.Printf("Created todo; %+v", r.URL.Path)
		components.Todo(todo).Render(r.Context(), w)
	})

	http.HandleFunc("DELETE /todos/{id}", func(w http.ResponseWriter, r *http.Request) {
		for idx, todo := range todos {
			if strconv.Itoa(todo.ID) == r.PathValue("id") {
				log.Printf("Deleting todo: %+v", todo)
				todos = append(todos[:idx], todos[idx+1:]...)
			}
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
