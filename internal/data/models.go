package data

import (
	"gorm.io/gorm"
	"time"
)

type Page struct {
	ID              string         `json:"id" gorm:"primarykey;type:varchar(36)"`
	Title           string         `json:"title"`
	Content         string         `json:"content" gorm:"type:text"`
	ParentID        string         `json:"parent_id" gorm:"type:varchar(36)"`
	BackgroundColor string         `json:"background_color" gorm:"type:varchar(20)"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Type            string         `json:"type"`
	Color           string         `json:"color"`
}

// 添加表名方法
func (Page) TableName() string {
	return "pages"
}

type Block struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"` // text, image, code, etc.
	Content   string    `json:"content"`
	PageID    string    `json:"page_id"`
	Order     int       `json:"order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
