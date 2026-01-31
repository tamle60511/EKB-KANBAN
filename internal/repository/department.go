package repository

import (
	"context"
	"cqs-kanban/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type (
	departmentRepo struct {
		db *gorm.DB
	}
	DepartmentRepo interface {
		Create(ctx context.Context, input models.Department) error
		GetByID(ctx context.Context, id uint) (*models.Department, error)
		GetAll(ctx context.Context) ([]models.Department, error)
		Update(ctx context.Context, id uint, input models.Department) error
		Delete(ctx context.Context, id uint) error
		Count(ctx context.Context) (int64, error)
	}
)

func NewDepartmentRepo(db *gorm.DB) DepartmentRepo {
	return &departmentRepo{
		db: db,
	}
}

func (d *departmentRepo) Create(ctx context.Context, input models.Department) error {
	if err := d.db.WithContext(ctx).Create(&input).Error; err != nil {
		return fmt.Errorf("failed to create department %w", err)
	}
	return nil
}

func (d *departmentRepo) GetByID(ctx context.Context, id uint) (*models.Department, error) {
	var department models.Department
	if err := d.db.WithContext(ctx).First(&department, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get department by id %w", err)
	}
	return &department, nil
}
func (d *departmentRepo) GetAll(ctx context.Context) ([]models.Department, error) {
	var departments []models.Department
	if err := d.db.WithContext(ctx).Find(&departments).Error; err != nil {
		return nil, fmt.Errorf("failed to get all departments %w", err)
	}
	return departments, nil
}

func (d *departmentRepo) Update(ctx context.Context, id uint, input models.Department) error {
	if err := d.db.WithContext(ctx).Model(&models.Department{}).Where("id = ?", id).Updates(input).Error; err != nil {
		return fmt.Errorf("failed to update department %w", err)
	}
	return nil
}

func (d *departmentRepo) Delete(ctx context.Context, id uint) error {
	if err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Department{}).Error; err != nil {
		return fmt.Errorf("failed to delete department %w", err)
	}
	return nil
}
func (d *departmentRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := d.db.Model(&models.Department{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count departments %w", err)
	}
	return count, nil
}
