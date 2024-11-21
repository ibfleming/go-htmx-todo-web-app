package handlers

import (
	"net/http"
	"strconv"
	"zion/internal/middleware/auth"
	"zion/internal/storage"
	"zion/internal/storage/schema"
	"zion/templates"
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

	// 3. Create todo
	todo, err := h.todos.CreateTodo(schema.Todo{
		UserID:      user.ID,
		Title:       title,
		Description: desc,
	})

	// 4. Check if no issue with creating todo
	if err != nil {
		http.Error(w, "failed to create todo", http.StatusInternalServerError)
		return
	}

	// 5. Render template
	err = templates.SingleTodo(todo).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}

func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// 1. Get parameters
	id := r.PathValue("id")
	user := auth.GetUser(r.Context())

	// 2. Delete todo
	err := h.todos.DeleteTodo(id, user.ID)
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

func (h *TodoHandler) EditItem(w http.ResponseWriter, r *http.Request) {
	// 1. Get parameters
	id := r.PathValue("id")

	item, err := h.todos.GetTodoItemByID(id)
	if err != nil {
		http.Error(w, "failed to get todo item", http.StatusInternalServerError)
		return
	}

	err = templates.EditTodoItem(item).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}

func (h *TodoHandler) UpdateItemContent(w http.ResponseWriter, r *http.Request) {
	// 1. Get parameters
	id := r.PathValue("id")
	content := r.FormValue("content")

	// 2. Update todo item content
	item, err := h.todos.UpdateTodoItemContent(id, content)
	if err != nil {
		http.Error(w, "failed to update todo item", http.StatusInternalServerError)
		return
	}

	// 3. Render template
	err = templates.NotCompletedTodoItem(item).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}

func (h *TodoHandler) ToggleItemCheck(w http.ResponseWriter, r *http.Request) {
	var checked bool

	// 1. Get parameters
	id := r.PathValue("id")
	check := r.FormValue("checked")
	r.ParseForm()

	// 2. Update todo item checked
	if check == "on" {
		checked = true
		item, err := h.todos.UpdateTodoItemChecked(id, checked)
		if err != nil {
			http.Error(w, "failed to update todo item", http.StatusInternalServerError)
			return
		}
		// 3. Render template
		err = templates.CompletedTodoItem(item).Render(r.Context(), w)
		if err != nil {
			http.Error(w, "failed to render template", http.StatusInternalServerError)
			return
		}
	} else {
		checked = false
		item, err := h.todos.UpdateTodoItemChecked(id, checked)
		if err != nil {
			http.Error(w, "failed to update todo item", http.StatusInternalServerError)
			return
		}
		// 3. Render template
		err = templates.NotCompletedTodoItem(item).Render(r.Context(), w)
		if err != nil {
			http.Error(w, "failed to render template", http.StatusInternalServerError)
			return
		}
	}
}

func (h *TodoHandler) AddItem(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	item, err := h.todos.AddTodoItemToTodo(&schema.TodoItem{
		TodoID:  uint(idUint),
		Content: "",
		Checked: false,
	})

	if err != nil {
		http.Error(w, "failed to add item", http.StatusInternalServerError)
		return
	}

	err = templates.SingleTodoItem(*item).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}

func (h *TodoHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	todoId := r.PathValue("todoId")
	itemId := r.PathValue("itemId")

	err := h.todos.DeleteTodoItemByID(itemId)
	if err != nil {
		http.Error(w, "failed to delete todo item", http.StatusInternalServerError)
		return
	}

	todoIdUint, err := strconv.ParseUint(todoId, 10, 32)
	if err != nil {
		http.Error(w, "invalid todoId", http.StatusBadRequest)
		return
	}

	length, err := h.todos.GetTodoItemLenthByID(uint(todoIdUint))
	if err != nil {
		http.Error(w, "failed to get todo length", http.StatusInternalServerError)
		return
	}

	if length == 0 {
		err = templates.EmptyTodoItemList(todoId).Render(r.Context(), w)
		if err != nil {
			http.Error(w, "failed to render template", http.StatusInternalServerError)
			return
		}
		return
	}
}
