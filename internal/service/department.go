package service

import (
	"context"
	"cqs-kanban/internal/dto"
	"cqs-kanban/internal/models"
	"cqs-kanban/internal/repository"
)

type (
	departmentService struct {
		departmentRepo repository.DepartmentRepo
	}
	DepartmentService interface {
		Create(ctx context.Context, req dto.DepartmenCreateReq) error
		Update(ctx context.Context, id uint, input dto.DepartmentUpdateReq) error
		GetByID(ctx context.Context, id uint) (*dto.DepartmentRes, error)
		GetAll(ctx context.Context) ([]dto.DepartmentRes, error)
		Delete(ctx context.Context, id uint) error
		Count(ctx context.Context) (int64, error)
	}
)

func NewDepartmentService(departmentRepo repository.DepartmentRepo) DepartmentService {
	return &departmentService{
		departmentRepo: departmentRepo,
	}
}

func (d *departmentService) Create(ctx context.Context, req dto.DepartmenCreateReq) error {
	return d.departmentRepo.Create(ctx, models.Department{
		Name: req.Name,
		Code: req.Code,
		Desc: req.Desc,
	})
}

func (d *departmentService) Update(ctx context.Context, id uint, input dto.DepartmentUpdateReq) error {
	dept := models.Department{}
	if input.Name != nil {
		dept.Name = *input.Name
	}
	if input.Code != nil {
		dept.Code = *input.Code
	}
	if input.Desc != nil {
		dept.Desc = *input.Desc
	}
	return d.departmentRepo.Update(ctx, id, dept)
}

func (d *departmentService) GetByID(ctx context.Context, id uint) (*dto.DepartmentRes, error) {
	dept, err := d.departmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if dept == nil {
		return nil, nil
	}
	return &dto.DepartmentRes{
		Name: dept.Name,
		Code: dept.Code,
		Desc: dept.Desc,
	}, nil
}

func (d *departmentService) GetAll(ctx context.Context) ([]dto.DepartmentRes, error) {
	depts, err := d.departmentRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var res []dto.DepartmentRes
	for _, dept := range depts {
		res = append(res, dto.DepartmentRes{
			ID:        dept.ID,
			Name:      dept.Name,
			Code:      dept.Code,
			Desc:      dept.Desc,
			CreatedAt: dept.CreatedAt.Format("2006-01-02"),
		})
	}
	return res, nil
}

func (d *departmentService) Delete(ctx context.Context, id uint) error {
	return d.departmentRepo.Delete(ctx, id)
}
func (d *departmentService) Count(ctx context.Context) (int64, error) {
	return d.departmentRepo.Count(ctx)
}
