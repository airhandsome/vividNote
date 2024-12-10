package service

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"path/filepath"
	"vividnote/internal/data"
)

func InitDB() (*gorm.DB, error) {
	dbPath := filepath.Join(".", "vividnote.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移数据库结构

	err = db.AutoMigrate(&data.Page{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
