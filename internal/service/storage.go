package service

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"vividnote/internal/data"
)

type StorageService struct {
	db *gorm.DB
}

func NewStorageService(db *gorm.DB) *StorageService {
	return &StorageService{db: db}
}

// 创建页面
func (s *StorageService) CreatePage(page *data.Page) error {
	fmt.Println(page)
	if page.ID == "" {
		page.ID = uuid.New().String()
	}
	return s.db.Create(page).Error
}

// 获取页面
func (s *StorageService) GetPageByID(id string) (*data.Page, error) {
	var page data.Page
	err := s.db.First(&page, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// 获取子页面
func (s *StorageService) GetChildPages(parentID string) ([]data.Page, error) {
	var pages []data.Page
	err := s.db.Where("parent_id = ?", parentID).Find(&pages).Error
	return pages, err
}

// 更新页面
func (s *StorageService) UpdatePage(page *data.Page) error {
	// 使用 Select 指定要更新的字段，排除 created_at
	result := s.db.Model(&data.Page{}).
		Select("title", "content", "color", "parent_id", "updated_at", "background_color").
		Where("id = ?", page.ID).
		Updates(page)

	if result.Error != nil {
		return fmt.Errorf("更新页面失败: %v", result.Error)
	}
	return nil
}

func (s *StorageService) DeletePage(pageID string) {
	s.db.Delete(&data.Page{}, "id = ?", pageID)
}
