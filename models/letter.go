package models

import "time"

type Letter struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `json:"user_id"`
	TypeID      uint       `json:"type_id"`
	Status      string     `gorm:"type:enum('pending','accepted','rejected');default:'pending'" json:"status"`
	RejectReason string    `json:"reject_reason"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	User        User       `gorm:"foreignKey:UserID"`
	LetterType  LetterType `gorm:"foreignKey:TypeID"`
}
