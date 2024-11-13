package handlers

import (
	"log"
	"net/http"
	"zion/internal/middleware/auth"
	"zion/internal/storage"
	"zion/internal/storage/db"
	"zion/templates"

	"github.com/go-chi/chi/v5"
)

type TodoHandler struct {
	users storage.UserStorageInterface
	todos storage.TodoStorageInterface
}

type TodoHandlerParams struct {
	Users storage.UserStorageInterface
	Todos storage.TodoStorageInterface
}

func NewTodoHandler(params TodoHandlerParams) *TodoHandler {
	return &TodoHandler{
		users: params.Users,
		todos: params.Todos,
	}
}

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	// 1. Get parameters
	title := r.FormValue("title")
	desc := r.FormValue("description")
	user := auth.GetUser(r.Context())

	// 2. Check input parameters
	if user == nil {
		http.Error(w, "no user current user found", http.StatusUnauthorized)
		return
	}

	if title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	// 3. Check if no todos found prior to adding this so we can remove empty message
	todos, _ := h.todos.GetTodosByUserID(user.ID)
	if len(todos) == 0 {
		log.Print("No todos found")
		w.Header().Set("HX-Trigger", "removeEmptyMessage")
	}

	// 4. Create todo
	todo, err := h.todos.CreateTodo(db.Todo{
		UserID:      user.ID,
		Title:       title,
		Description: desc,
	})

	// 5. Check if no issue with creating todo
	if err != nil {
		http.Error(w, "failed to create todo", http.StatusInternalServerError)
		return
	}

	// 6. Render template
	err = templates.SingleTodo(todo).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}

func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// 1. Get parameters
	user := auth.GetUser(r.Context())

	// 2. Delete todo
	err := h.todos.DeleteTodo(chi.URLParam(r, "id"), user.ID)
	if err != nil {
		http.Error(w, "failed to delete todo", http.StatusInternalServerError)
		return
	}

	// 3. Get todos
	todos, err := h.todos.GetTodosByUserID(user.ID)
	if err != nil {
		http.Error(w, "failed to get todos", http.StatusInternalServerError)
		return
	}

	// 4. Optionally render empty todo list
	if len(todos) == 0 {
		err := templates.EmptyTodoList().Render(r.Context(), w)
		if err != nil {
			http.Error(w, "failed to render template", http.StatusInternalServerError)
			return
		}
	}
}

func (h *TodoHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	// 1. Get parameters
	user := auth.GetUser(r.Context())

	// 2. Delete todo
	err := h.todos.DeleteAllTodos(user.ID)
	if err != nil {
		http.Error(w, "failed to delete all todos", http.StatusInternalServerError)
		return
	}

	err = templates.EmptyTodoList().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}

func (h *TodoHandler) List(w http.ResponseWriter, r *http.Request) {
	// 1. Get parameters
	user := auth.GetUser(r.Context())

	// 2. Check input parameters
	if user == nil {
		http.Error(w, "no user current user found", http.StatusUnauthorized)
		return
	}

	// 3. Get todos
	todos, err := h.todos.GetTodosByUserID(user.ID)
	if err != nil {
		http.Error(w, "failed to get todos", http.StatusInternalServerError)
		return
	}

	// 4. Render template

	// 4a Render empty todo list
	if len(todos) == 0 {
		err = templates.EmptyTodoList().Render(r.Context(), w)
		if err != nil {
			http.Error(w, "failed to render template", http.StatusInternalServerError)
			return
		}
		return
	}

	// 4b Render todo list
	err = templates.TodoList(todos).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}
