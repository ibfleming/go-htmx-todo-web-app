package db

import (
	"log"
	zerr "zion/internal/errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	SessionID string `gorm:"unique"`
	UserID    uint
	User      User `gorm:"foreignKey:UserID"`
}

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Todos    []Todo `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

type Todo struct {
	gorm.Model
	UserID      uint       `gorm:"not null;index"`
	Title       string     `gorm:"not null;size:255"`
	Description string     `gorm:"default:null;size:255"`
	Items       []TodoItem `gorm:"foreignKey:TodoID;constraint:OnDelete:CASCADE;"`
}

type TodoItem struct {
	gorm.Model
	TodoID  uint   `gorm:"not null;index"`
	Content string `gorm:"not null;size:500"`
	Checked bool   `gorm:"default:false;index"`
}

type UserStorage interface {
	CreateUser(email, password string) error
	GetUser(email string) (*User, error)
}

type SessionStorage interface {
	CreateSession(sessionID string, userID uint) error
	GetUserFromSession(sessionID, userId string) (*User, error)
}

var models = []interface{}{&User{}, &Session{}, &Todo{}}

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
	log.Println("Migrating models into database...")
	log.Println("Models:")
	for _, model := range models {
		log.Printf("\tModel: %T", model)
	}
	return db.AutoMigrate(models...)
}
