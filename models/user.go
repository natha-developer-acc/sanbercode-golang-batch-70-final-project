package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RoleID    uint      `json:"role_id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"unique" json:"email"`
	Password  string    `json:"-"` // jangan expose password
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Role      Role      `gorm:"foreignKey:RoleID"`
}

// untuk input register
type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
