package model

import "sync"

type Todo struct {
	ID    int
	Title string
}

type TodoDB interface {
	AddTodo(title string) Todo
	ListTodos() []Todo
	RemoveTodo(ID int) (Todo, bool)
}

var _ TodoDB = &todoDB{}

type todoDB struct {
	sync.Mutex
	nextID int
	todos  []Todo
}

func NewTodoDB() TodoDB {
	return &todoDB{
		todos: []Todo{
			{ID: 1, Title: "Clean the dishes"},
			{ID: 2, Title: "Learn TailwindCSS"},
			{ID: 3, Title: "Workout"},
		}, nextID: 4,
	}
}

// AddTodo implements TodoDB.
func (t *todoDB) AddTodo(title string) Todo {
	t.Lock()
	defer t.Unlock()

	newTodo := Todo{
		ID:    t.nextID,
		Title: title,
	}
	t.todos = append(t.todos, newTodo)

	t.nextID++
	return newTodo
}

// ListTodos implements TodoDB.
func (t *todoDB) ListTodos() []Todo {
	t.Lock()
	defer t.Unlock()

	res := make([]Todo, 0, len(t.todos))
	res = append(res, t.todos...)

	return res
}

// RemoveTodo implements TodoDB.
func (t *todoDB) RemoveTodo(ID int) (Todo, bool) {
	t.Lock()
	defer t.Unlock()

	for idx, todo := range t.todos {
		if ID == todo.ID {
			removed := todo
			t.todos = append(t.todos[:idx], t.todos[idx+1:]...)
			return removed, true
		}
	}

	return Todo{}, false
}
