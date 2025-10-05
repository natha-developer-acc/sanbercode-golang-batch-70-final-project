package models

import "time"

type Setting struct {
    ID            uint      `gorm:"primaryKey" json:"id"`
    UserID        uint      `json:"user_id"`
    TelegramChatID string   `gorm:"column:telegram_chatid" json:"telegram_chatid"`
    WANumber      string    `json:"wa_number"`
    AllowTelegram string    `gorm:"type:enum('yes','no');default:'no'" json:"allow_telegram"`
    AllowWA       string    `gorm:"type:enum('yes','no');default:'no'" json:"allow_wa"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
    User          User      `gorm:"foreignKey:UserID"`
}
