Here is a step-by-step guide to implementing an efficient TODO application with Go, `chi` router, Templ, htmx, and GORM for PostgreSQL. Each step is broken down to facilitate a seamless implementation, from setting up models to integrating htmx for adding and deleting TODOs.

---

### Step 1: Define Models with GORM

Define the models for `User` and `Todo` to represent users and their TODO items.

```go
// models.go
type User struct {
    ID       uint   `gorm:"primaryKey"`
    Username string `gorm:"uniqueIndex"`
    Todos    []Todo
}

type Todo struct {
    ID     uint   `gorm:"primaryKey"`
    UserID uint   `gorm:"index"`
    Title  string
    Done   bool
}
```

### Step 2: Create Database Functions

These helper functions will perform CRUD operations for TODOs.

```go
// db.go
func GetTodosForUser(db *gorm.DB, userID uint) ([]Todo, error) {
    var todos []Todo
    err := db.Where("user_id = ?", userID).Find(&todos).Error
    return todos, err
}

func CreateTodoForUser(db *gorm.DB, userID uint, title string) (Todo, error) {
    todo := Todo{UserID: userID, Title: title}
    err := db.Create(&todo).Error
    return todo, err
}

func DeleteTodoByID(db *gorm.DB, todoID uint) error {
    return db.Delete(&Todo{}, todoID).Error
}
```

### Step 3: Define Templ Components

1. **SingleTodoItem**: A component to render a single TODO item.
2. **TodoList**: Renders a list of TODO items.
3. **TodoPage**: Main page with the list and `Add New Todo` button.

```go
// templates.templ
templ SingleTodoItem(todo Todo) {
    <li id={"todo_" + fmt.Sprint(todo.ID)}>
        <span class={ "todo-item", templ.KV("done", todo.Done) }>{ todo.Title }</span>
        <button hx-delete={"/todos/delete/" + fmt.Sprint(todo.ID)}
                hx-target={"#todo_" + fmt.Sprint(todo.ID)}
                hx-swap="outerHTML">
            Delete
        </button>
    </li>
}

templ TodoList(todos []Todo) {
    <ul id="todoList" class="todo-list">
        for _, todo := range todos {
            @SingleTodoItem(todo)
        }
    </ul>
}

templ TodoPage(todos []Todo) {
    <div class="todo-container">
        <h2>Your Todos</h2>
        @TodoList(todos)
        <form hx-post="/todos/create" hx-target="#todoList" hx-swap="beforeend">
            <input type="text" name="title" placeholder="New Todo" required/>
            <button type="submit" class="add-todo-button">Add</button>
        </form>
    </div>
}
```

### Step 4: Set Up Authentication Middleware

Add a middleware to extract the user ID from the session or token, making it accessible to handlers via context.

```go
// middleware.go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userID, err := GetUserIDFromSession(r) // Implement session extraction logic
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        ctx := context.WithValue(r.Context(), "userID", userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Step 5: Define Handlers for CRUD Operations

#### Handler to Render Todos List
Retrieves and renders the user’s TODO list with `TodoList` component.

```go
// handlers.go
func getTodosHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := r.Context().Value("userID").(uint)
        todos, _ := GetTodosForUser(db, userID)
        TodoList(todos).Render(r.Context(), w)
    }
}
```

#### Handler to Create a New Todo
Creates a TODO item and renders only the new item.

```go
func createTodoHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := r.Context().Value("userID").(uint)
        title := r.PostFormValue("title")

        newTodo, _ := CreateTodoForUser(db, userID, title)

        // Render only the new TODO item
        SingleTodoItem(newTodo).Render(r.Context(), w)
    }
}
```

#### Handler to Delete a Todo
Deletes the specified TODO and re-renders the list (if needed).

```go
func deleteTodoHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        todoID, _ := strconv.Atoi(chi.URLParam(r, "id"))
        _ = DeleteTodoByID(db, uint(todoID))
        http.Error(w, "", http.StatusNoContent) // Indicate deletion with a 204 response
    }
}
```

### Step 6: Set Up Routes with `chi`

Configure `chi` routes, attaching the handlers and the authentication middleware.

```go
// router.go
func setupRoutes(r chi.Router, db *gorm.DB) {
    r.Route("/todos", func(r chi.Router) {
        r.Use(AuthMiddleware) // Applies middleware to all /todos routes
        r.Get("/", getTodosHandler(db))
        r.Post("/create", createTodoHandler(db))
        r.Delete("/delete/{id}", deleteTodoHandler(db))
    })
}
```

### Step 7: Serve the `TodoPage` Component on Login

When a user logs in, serve the `TodoPage` component, passing in their TODO list.

```go
func loginHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := r.Context().Value("userID").(uint) // Retrieved via AuthMiddleware
        todos, _ := GetTodosForUser(db, userID)
        TodoPage(todos).Render(r.Context(), w)
    }
}
```

### Step 8: Integrate Tailwind CSS (Optional)

Set up Tailwind for styling the UI. Create a `tailwind.css` with basic styling for buttons and TODO items.

```css
/* tailwind.css */
@tailwind base;
@tailwind components;
@tailwind utilities;

/* Custom styles */
.todo-container {
    @apply p-4 bg-gray-100 rounded-lg;
}

.add-todo-button {
    @apply bg-blue-500 text-white px-4 py-2 rounded mt-2;
}

.todo-item {
    @apply text-lg;
}
```

Generate your CSS for production:

```bash
npx tailwindcss -i ./tailwind.css -o ./tailwind_mini.css --minify
```

---

### Summary

1. **Models**: Define `User` and `Todo` models with GORM.
2. **Database Functions**: Implement `GetTodosForUser`, `CreateTodoForUser`, and `DeleteTodoByID`.
3. **Templ Components**: Create `SingleTodoItem`, `TodoList`, and `TodoPage`.
4. **Middleware**: Set up authentication to extract `userID`.
5. **Handlers**: Implement CRUD operations with optimized HTMX responses.
6. **Routes**: Configure `chi` routes.
7. **Login Handler**: Serve `TodoPage` on login.
8. **Tailwind CSS**: Add basic styling (optional).

This ordered guide provides a structured path to implementing a responsive and efficient TODO application. Let me know if you need more detail on any step!
