package db

import (
	"log"
	zerr "zion/internal/errors"
	"zion/internal/hash"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
}

type Session struct {
	gorm.Model
	SessionID string `gorm:"unique"`
	UserID    uint
	User      User `gorm:"foreignKey:UserID"`
}

type Todo struct {
	gorm.Model
	Title       string
	Description string
	Items       []TodoItem `gorm:"foreignKey:TodoID;constraint:OnDelete:CASCADE;"`
	UserID      uint
	User        User `gorm:"foreignKey:UserID"`
}

type TodoItem struct {
	gorm.Model
	Completed   bool
	Description string
	TodoID      uint
}

type UserStorage interface {
	CreateUser(email, password string) error
	GetUser(email string) (*User, error)
}

type SessionStorage interface {
	CreateSession(sessionID string, userID uint) error
	GetUserFromSession(sessionID, userId string) (*User, error)
}

var models = []interface{}{&User{}, &Session{}, &Todo{}, &TodoItem{}}

func Connect(url string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	return db, err
}

func CreateModels(db *gorm.DB, mode string, hash *hash.PasswordHash) error {
	switch mode {
	case "clear":
		{
			log.Println("clearing database...")
			err := db.Migrator().DropTable(models...)
			if err != nil {
				return zerr.ErrDropTables
			}
			return migrate(db)
		}
	case "seed":
		{
			log.Println("clearing database...")
			err := db.Migrator().DropTable(models...)
			if err != nil {
				return zerr.ErrDropTables
			}

			err = migrate(db)
			if err != nil {
				return err
			}

			log.Println("seeding database...")

			password, _ := hash.GenerateFromPassword("2211")

			mockUser := &User{
				Email:    "admin@localhost",
				Password: password,
			}

			db.Create(mockUser)

			mockTodo1 := &Todo{
				Title:       "Todo 1",
				Description: "This is the first todo",
				UserID:      mockUser.ID,
				User:        *mockUser,
			}

			db.Create(mockTodo1)

			mockTodoItem1 := &TodoItem{
				Completed:   false,
				Description: "This is the first todo item",
				TodoID:      mockTodo1.ID,
			}

			mockTodoItem2 := &TodoItem{
				Completed:   true,
				Description: "This is the second todo item",
				TodoID:      mockTodo1.ID,
			}

			db.Create(mockTodoItem1)
			db.Create(mockTodoItem2)

			db.Model(&mockTodo1).Association("Items").Append(mockTodoItem1)
			db.Model(&mockTodo1).Association("Items").Append(mockTodoItem2)

			db.Create(&Todo{
				Title:       "Todo 2",
				Description: "This is the second todo",
				UserID:      mockUser.ID,
				User:        *mockUser,
			})

			return err
		}
	default:
		{
			return migrate(db)
		}
	}
}

func migrate(db *gorm.DB) error {
	log.Println("migrating database...")
	log.Print("models:")
	for _, model := range models {
		log.Printf("%T", model)
	}
	return db.AutoMigrate(models...)
}
