package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"vividnote/internal/data"
	"vividnote/internal/service"
)

// App struct
type App struct {
	ctx     context.Context
	storage *service.StorageService
	db      *gorm.DB
}

// NewApp creates a new App application struct
func NewApp() *App {
	db, err := service.InitDB()
	if err != nil {
		panic(err)
	}

	return &App{
		storage: service.NewStorageService(db),
		db:      db,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
//func (a *App) Greet(name string) string {
//	return fmt.Sprintf("Hello %s, It's show time!", name)
//}

type Person struct {
	Name    string   `json:"name"`
	Age     uint8    `json:"age"`
	Address *Address `json:"address"`
}

type Address struct {
	Street   string `json:"street"`
	Postcode string `json:"postcode"`
}

func (a *App) Greet(p Person) string {
	return fmt.Sprintf("Hello %s (Age: %d)!", p.Name, p.Age)
}

// 创建新页面
func (a *App) CreatePage(page *data.Page) (*data.Page, error) {
	if page.ID == "" {
		page.ID = uuid.New().String()
	}
	// 确保设置了页面类型
	if page.Type == "" {
		page.Type = "original" // 默认为原始笔记
	}
	err := a.storage.CreatePage(page)
	if err != nil {
		return nil, err
	}
	return page, nil
}

// 更新页面
func (a *App) UpdatePage(page *data.Page) error {
	return a.storage.UpdatePage(page)
}

// 获取页面
func (a *App) GetPage(id string) (*data.Page, error) {
	return a.storage.GetPageByID(id)
}

// 获取子页面
func (a *App) GetChildPages(category string) ([]data.Page, error) {
	fmt.Printf("获取类别为 %s 的页面\n", category)
	var pages []data.Page
	result := a.db.Where("type = ?", category).Find(&pages)
	if result.Error != nil {
		return nil, fmt.Errorf("获取页面失败: %v", result.Error)
	}

	return pages, nil
}
func (a *App) DeletePage(pageID string) error {
	// 在数据库中删除页面的逻辑
	// 例如：db.Delete(&Page{}, "id = ?", pageID)
	if pageID == "" {
		return fmt.Errorf("pageID is required")
	}
	a.storage.DeletePage(pageID)
	fmt.Printf("删除页面: %s\n", pageID)
	return nil
}
