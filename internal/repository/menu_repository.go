package repository

import (
	"context"
	"cqs-kanban/internal/dto"
	"cqs-kanban/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type (
	menuRepo struct {
		db *gorm.DB
	}
	MenuRepo interface {
		Create(ctx context.Context, input dto.MenuCreateReq) error
		GetByDepartment(ctx context.Context, departmentID int64) ([]models.Menu, error)
		GetByDepartments(ctx context.Context, departmentIDs []int64) ([]models.Menu, error)
		GetByID(ctx context.Context, id int64) (*models.Menu, error)
		Update(ctx context.Context, id int64, input dto.MenuUpdateReq) error
		Delete(ctx context.Context, id int64) error
		GetAll(ctx context.Context) ([]models.Menu, error)
		GetItemsMenuIDs(ctx context.Context, id int64) ([]dto.MenuDetailItem, error)
	}
)

func NewMenuRepo(db *gorm.DB) MenuRepo {
	return &menuRepo{
		db: db,
	}
}

func (m *menuRepo) Create(ctx context.Context, input dto.MenuCreateReq) error {
	// Bắt đầu transaction
	tx := m.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Tạo menu parent
	menu := models.Menu{
		Title:        input.Detail.Title,
		Code:         input.Detail.Code,
		Route:        input.Detail.Route,
		Icon:         input.Detail.Icon,
		DepartmentID: input.DepartmentID,
		Level:        1,
		ParentID:     0,
	}

	if err := tx.Create(&menu).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create menu: %w", err)
	}

	// Tạo menu items
	if len(input.List) > 0 {
		var menuList []models.Menu
		for _, item := range input.List {
			menuList = append(menuList, models.Menu{
				Title:        item.Title,
				Code:         item.Code,
				Route:        item.Route,
				ReportID:     item.ReportID,
				Level:        2,
				ParentID:     menu.ID,
				DepartmentID: input.DepartmentID, // Thêm DepartmentID
			})
		}

		if err := tx.Create(&menuList).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create menu items: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
func (m *menuRepo) GetByDepartment(ctx context.Context, departmentID int64) ([]models.Menu, error) {
	detail, err := gorm.G[models.Menu](m.db).Where("department_id = ?", departmentID).Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("not found departmentID %w", err)
	}
	return detail, nil
}

func (m *menuRepo) GetByDepartments(ctx context.Context, departmentIDs []int64) ([]models.Menu, error) {
	detail, err := gorm.G[models.Menu](m.db).Where("department_id in ?", departmentIDs).Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("not found departmentID %w", err)
	}
	return detail, nil
}
func (m *menuRepo) GetItemsMenuIDs(ctx context.Context, id int64) ([]dto.MenuDetailItem, error) {
	menus, err := gorm.G[models.Menu](m.db).Where("parent_id = ?", id).Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu items %w", err)
	}

	items := make([]dto.MenuDetailItem, 0)

	for _, menu := range menus {
		items = append(items, dto.MenuDetailItem{
			ID:       menu.ID,
			Title:    menu.Title,
			Code:     menu.Code,
			Route:    menu.Route,
			ReportID: menu.ReportID,
		})
	}
	return items, nil
}
func (m *menuRepo) GetByID(ctx context.Context, id int64) (*models.Menu, error) {
	menus, err := gorm.G[models.Menu](m.db).Where("id = ?", id).Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu %w", err)
	}

	if len(menus) < 1 {
		return nil, fmt.Errorf("not found menu")
	}
	return &menus[0], nil
}
func (m *menuRepo) Update(ctx context.Context, id int64, input dto.MenuUpdateReq) error {
	menu, err := m.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get menu %w", err)
	}

	if input.Detail != nil {
		menu.Title = input.Detail.Title
		menu.Code = input.Detail.Code
		menu.Route = input.Detail.Route
	}

	if input.DepartmentID != nil {
		menu.DepartmentID = *input.DepartmentID
	}

	if err := m.db.WithContext(ctx).Save(menu).Error; err != nil {
		return fmt.Errorf("failed to update menu %w", err)
	}

	if input.List != nil {
		if err := m.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.WithContext(ctx).Where("parent_id = ?", id).Delete(&models.Menu{}).Error; err != nil {
				return fmt.Errorf("failed to delete existing sub-menus %w", err)
			}

			var menuList []models.Menu
			for _, item := range *input.List {
				// ✅ FIX: Kiểm tra nil trước khi dereference
				if item.Title == nil || item.Code == nil || item.Route == nil || item.ReportID == nil {
					return fmt.Errorf("sub-menu fields cannot be null")
				}

				menuList = append(menuList, models.Menu{
					Title:    *item.Title,
					Code:     *item.Code,
					Route:    *item.Route,
					ReportID: *item.ReportID,
					Level:    2,
					ParentID: menu.ID,
				})
			}

			if len(menuList) > 0 {
				if err := tx.WithContext(ctx).Create(&menuList).Error; err != nil {
					return fmt.Errorf("failed to create new sub-menus %w", err)
				}
			}

			return nil
		}); err != nil {
			return fmt.Errorf("failed to update sub-menus %w", err)
		}
	}

	return nil
}
func (m *menuRepo) Delete(ctx context.Context, id int64) error {
	if err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Where("parent_id = ?", id).Delete(&models.Menu{}).Error; err != nil {
			return fmt.Errorf("failed to delete sub-menus %w", err)
		}
		if err := tx.WithContext(ctx).Where("id = ?", id).Delete(&models.Menu{}).Error; err != nil {
			return fmt.Errorf("failed to delete menu %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to delete menus %w", err)
	}
	return nil
}

func (m *menuRepo) GetAll(ctx context.Context) ([]models.Menu, error) {
	menus, err := gorm.G[models.Menu](m.db).Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get menus %w", err)
	}
	return menus, nil
}
