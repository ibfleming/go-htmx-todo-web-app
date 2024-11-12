package handlers

import (
	"net/http"
	"zion/internal/middleware/auth"
	"zion/internal/storage"
	"zion/internal/storage/db"
	"zion/templates"
)

type Index struct {
	todos storage.TodoStorageInterface
}

type IndexParams struct {
	Todos storage.TodoStorageInterface
}

func NewGetIndexHandler(params IndexParams) *Index {
	return &Index{
		todos: params.Todos,
	}
}

func (h *Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// (1) Get User from context if it exists
	user := auth.GetUser(r.Context())

	// (2) Fetch all todos
	var todos []db.Todo
	if user != nil {
		var err error
		todos, err = h.todos.GetAllTodosForUser(user.ID)
		if err != nil {
			http.Error(w, "error getting todos", http.StatusInternalServerError)
		}
	}

	// (3) Render the index template
	err := templates.Index(user, todos).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}
}
