package db

import (
	"log"
	zerr "zion/internal/errors"
	"zion/internal/storage"
	s "zion/internal/storage/schema"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CreateModelsParams struct {
	DB    *gorm.DB
	Users storage.UserStorageInterface
	Todos storage.TodoStorageInterface
	Mode  string
}

var Models = []interface{}{&s.User{}, &s.Session{}, &s.Todo{}, &s.TodoItem{}}

func Connect(url string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	return db, err
}

func CreateModels(p CreateModelsParams) error {
	switch p.Mode {
	case "clear":
		{
			log.Println("clearing the store...")
			err := p.DB.Migrator().DropTable(Models...)
			if err != nil {
				return zerr.ErrDropTables
			}
			return migrate(p.DB)
		}
	case "seed":
		{
			log.Println("clearing the database...")
			err := p.DB.Migrator().DropTable(Models...)
			if err != nil {
				return zerr.ErrDropTables
			}
			log.Println("seeding the database...")

			err = migrate(p.DB)

			p.Users.CreateUser("admin@localhost", "2211")

			admin, _ := p.Users.GetUser("admin@localhost")

			p.Todos.CreateTodo(s.Todo{
				UserID:      admin.ID,
				Title:       "Mock todo",
				Description: "This is a mock todo",
				Items: []s.TodoItem{
					{
						Content: "Mock todo item",
						Checked: false,
					},
					{
						Content: "Mock todo item",
						Checked: true,
					},
				},
			})

			p.Todos.CreateTodo(s.Todo{
				UserID:      admin.ID,
				Title:       "Another mock todo",
				Description: "This is another mock todo",
				Items:       []s.TodoItem{},
			})

			return err
		}
	default:
		{
			return migrate(p.DB)
		}
	}
}

func migrate(db *gorm.DB) error {
	log.Println("adding models to database...")
	log.Println("models:")
	for _, model := range Models {
		log.Printf("\t%T", model)
	}
	return db.AutoMigrate(Models...)
}
