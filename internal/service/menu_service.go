package service

import (
	"context"
	"cqs-kanban/internal/dto"
	"cqs-kanban/internal/models"
	"cqs-kanban/internal/repository"
	"fmt"
	"strings"
)

type (
	menuService struct {
		menuRepo repository.MenuRepo
	}
	MenuService interface {
		Create(ctx context.Context, req dto.MenuCreateReq) error
		GetMenu(ctx context.Context, req dto.MenuDetailReq) (*dto.MenuDetailRes, error)
		Update(ctx context.Context, id int64, req dto.MenuUpdateReq) error
		Delete(ctx context.Context, id int64) error
		GetAll(ctx context.Context) ([]dto.MenuDetailRes, error)
		GetByID(ctx context.Context, id int64) (*dto.MenuDetailRes, error)
		GetItemsMenuIDs(ctx context.Context, id int64) ([]dto.MenuDetailItem, error)
	}
)

func NewMenuService(menuRepo repository.MenuRepo) MenuService {
	return &menuService{
		menuRepo: menuRepo,
	}
}
func (m *menuService) Create(ctx context.Context, req dto.MenuCreateReq) error {
	// Validate department exists
	if req.DepartmentID <= 0 {
		return fmt.Errorf("invalid department_id")
	}

	// Validate menu detail
	if strings.TrimSpace(req.Detail.Title) == "" {
		return fmt.Errorf("menu title is required")
	}
	if strings.TrimSpace(req.Detail.Code) == "" {
		return fmt.Errorf("menu code is required")
	}

	// Check duplicate codes
	codes := make(map[string]bool)
	parentCode := strings.ToUpper(strings.TrimSpace(req.Detail.Code))
	codes[parentCode] = true

	for _, item := range req.List {
		itemCode := strings.ToUpper(strings.TrimSpace(item.Code))
		if itemCode == "" {
			continue
		}
		if codes[itemCode] {
			return fmt.Errorf("duplicate menu code: %s", itemCode)
		}
		codes[itemCode] = true
	}

	return m.menuRepo.Create(ctx, req)
}

func (m *menuService) GetMenu(ctx context.Context, req dto.MenuDetailReq) (*dto.MenuDetailRes, error) {
	var menus []models.Menu
	var err error

	if req.DepartmentID != nil && *req.DepartmentID > 0 {
		menus, err = m.menuRepo.GetByDepartment(ctx, *req.DepartmentID)
		if err != nil {
			return nil, fmt.Errorf("not found department id: %w", err)
		}
	} else {
		menus, err = m.menuRepo.GetAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get all menus: %w", err)
		}
	}

	ms := make(map[int64][]models.Menu)
	for _, menu := range menus {
		if menu.Level != 1 && menu.ParentID != 0 {
			ms[menu.ParentID] = append(ms[menu.ParentID], menu)
		}
	}
	var menuDetails []dto.MenuDetail
	for _, menu := range menus {
		if menu.Level == 1 || menu.ParentID == 0 {

			menuDetail := dto.MenuDetail{
				ID:           menu.ID,
				Title:        menu.Title,
				Code:         menu.Code,
				Route:        menu.Route,
				Icon:         menu.Icon,
				DepartmentID: menu.DepartmentID,
				List:         make([]dto.MenuDetailItem, 0),
				CreatedAt:    menu.CreatedAt.Format("2006-01-02"),
				UpdatedAt:    menu.UpdatedAt.Format("2006-01-02"),
			}

			// Add children
			if children, exists := ms[menu.ID]; exists {
				for _, child := range children {
					item := dto.MenuDetailItem{
						ID:       child.ID,
						Title:    child.Title,
						Code:     child.Code,
						Route:    child.Route,
						ReportID: child.ReportID,
					}
					menuDetail.List = append(menuDetail.List, item)
				}
			}

			menuDetails = append(menuDetails, menuDetail)
		}
	}

	res := &dto.MenuDetailRes{
		Menu: menuDetails,
	}

	return res, nil
}
func (m *menuService) Update(ctx context.Context, id int64, req dto.MenuUpdateReq) error {
	return m.menuRepo.Update(ctx, id, req)
}
func (m *menuService) Delete(ctx context.Context, id int64) error {
	return m.menuRepo.Delete(ctx, id)
}
func (m *menuService) GetAll(ctx context.Context) ([]dto.MenuDetailRes, error) {
	menus, err := m.menuRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var res []dto.MenuDetailRes
	menuMap := make(map[int64][]dto.MenuDetailItem)

	for _, menu := range menus {
		if menu.ParentID != 0 {
			menuMap[menu.ParentID] = append(menuMap[menu.ParentID], dto.MenuDetailItem{
				ID:       menu.ID,
				Title:    menu.Title,
				Code:     menu.Code,
				Route:    menu.Route,
				Icon:     menu.Icon,
				ReportID: menu.ReportID,
			})
		}
	}
	for _, menu := range menus {
		if menu.ParentID == 0 {
			menuDetail := dto.MenuDetailRes{
				Menu: []dto.MenuDetail{
					{
						ID:           menu.ID,
						Title:        menu.Title,
						Code:         menu.Code,
						List:         menuMap[menu.ID],
						Route:        menu.Route,
						Icon:         menu.Icon,
						DepartmentID: menu.DepartmentID,
						CreatedAt:    menu.CreatedAt.Format("2006-01-02"),
						UpdatedAt:    menu.UpdatedAt.Format("2006-01-02"),
					},
				},
			}
			res = append(res, menuDetail)
		}
	}

	return res, nil
}
func (m *menuService) GetItemsMenuIDs(ctx context.Context, id int64) ([]dto.MenuDetailItem, error) {
	return m.menuRepo.GetItemsMenuIDs(ctx, id)
}
func (m *menuService) GetByID(ctx context.Context, id int64) (*dto.MenuDetailRes, error) {
	menu, err := m.menuRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if menu == nil {
		return nil, nil
	}
	menuItems, err := m.menuRepo.GetItemsMenuIDs(ctx, menu.ID)
	if err != nil {
		return nil, err
	}
	detail := &dto.MenuDetailRes{
		Menu: []dto.MenuDetail{
			{
				ID:           menu.ID,
				Title:        menu.Title,
				Code:         menu.Code,
				Route:        menu.Route,
				Icon:         menu.Icon,
				DepartmentID: menu.DepartmentID,
				List:         menuItems,
				CreatedAt:    menu.CreatedAt.Format("2006-01-02"),
				UpdatedAt:    menu.UpdatedAt.Format("2006-01-02"),
			},
		},
	}
	return detail, nil
}
