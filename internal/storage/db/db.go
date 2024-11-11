package db

import (
	"log"
	zerr "zion/internal/errors"

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

func CreateModels(db *gorm.DB, mode string) error {
	switch mode {
	case "clear":
		{
			log.Println("🔨 Clearing the store...")
			err := db.Migrator().DropTable(models...)
			if err != nil {
				return zerr.ErrDropTables
			}
			return migrate(db)
		}
	case "seed":
		{
			log.Println("🔨 Seeding the store...")
			return migrate(db)
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
