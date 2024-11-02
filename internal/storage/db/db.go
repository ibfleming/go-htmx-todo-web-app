package db

import (
	"log"
	"os"

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

type UserStorage interface {
	CreateUser(email, password string) error
	GetUser(email string) (*User, error)
}

type SessionStorage interface {
	CreateSession(sessionID string, userID uint) error
	GetUserFromSession(sessionID, userId string) (*User, error)
}

var models = []interface{}{&User{}, &Session{}}

func Connect(url string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	return db, err
}

func CreateModels(db *gorm.DB) {
	mode := os.Getenv("DATABASE_MODE")
	// Migrate tables to the store depending on mode
	switch mode {
	case "clear":
		{
			log.Println("🔨 Clearing the store...")
			err := db.Migrator().DropTable(models...)
			if err != nil {
				log.Fatalf("❌ Failed to drop the tables in the store: %v", err)
			}
			migrate(db)
		}
	case "seed":
		{
			log.Println("🔨 Seeding the store...")
			migrate(db)
		}
	default:
		{
			migrate(db)
		}
	}
}

func migrate(db *gorm.DB) {
	log.Println("📦 Migrating models:")
	for _, model := range models {
		log.Printf("\tModel: %T", model)
	}
	err := db.AutoMigrate(models...)
	if err != nil {
		log.Fatalf("❌ Failed to migrate models: %v", err)
	}
}
