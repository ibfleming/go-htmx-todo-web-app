package schema

import "gorm.io/gorm"

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
