package repository

import (
	"context"
	"cqs-kanban/internal/dto"
	"cqs-kanban/internal/models"
	"cqs-kanban/internal/utils"
	"fmt"

	"gorm.io/gorm"
)

type (
	userRepo struct {
		db *gorm.DB
	}
	UserRepo interface {
		Create(ctx context.Context, input dto.UserCreateReq) error
		GetByID(ctx context.Context, id int64) (*models.User, error)
		GetByUsername(ctx context.Context, username string) (*models.User, error)
		Update(ctx context.Context, id int64, input dto.UserUpdateReq) error
		UpdatePassword(ctx context.Context, id int64, password string) error
		Delete(ctx context.Context, id int64) error
		GetAll(ctx context.Context) ([]models.User, error)
		GetUser(ctx context.Context, req dto.UserDetailReq) (*dto.UserDetailRes, error)
		Count(ctx context.Context) (int64, error)
	}
)

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}

func (u *userRepo) Create(ctx context.Context, input dto.UserCreateReq) error {
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password %w", err)
	}

	user := &models.User{
		Username:     input.Username,
		FullName:     input.FullName,
		Email:        input.Email,
		DepartmentID: input.DepartmentID,
		Password:     hashedPassword,
	}

	if err := u.db.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user %w", err)
	}
	return nil
}

func (u *userRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	users, err := gorm.G[models.User](u.db).Where("id = ?", id).Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users %w", err)
	}

	if len(users) < 1 {
		return nil, fmt.Errorf("not found user")
	}
	return &users[0], nil
}

func (u *userRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	users, err := gorm.G[models.User](u.db).Where("username = ?", username).Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users %w", err)
	}

	if len(users) < 1 {
		return nil, fmt.Errorf("not found user")
	}
	return &users[0], nil
}

func (u *userRepo) Update(ctx context.Context, id int64, input dto.UserUpdateReq) error {
	updates := make(map[string]interface{})
	if input.Username != nil {
		updates["username"] = *input.Username
	}
	if input.FullName != nil {
		updates["full_name"] = *input.FullName
	}
	if input.Email != nil {
		updates["email"] = *input.Email
	}
	if input.DepartmentID != nil {
		updates["department_id"] = *input.DepartmentID
	}
	if input.Password != nil {
		hashedPassword, err := utils.HashPassword(*input.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password %w", err)
		}
		updates["password"] = hashedPassword
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	if err := u.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update user %w", err)
	}
	return nil
}
func (u *userRepo) UpdatePassword(ctx context.Context, id int64, password string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password %w", err)
	}
	if err := u.db.Model(&models.User{}).Where("id = ?", id).Update("password", hashedPassword).Error; err != nil {
		return fmt.Errorf("failed to update user password %w", err)
	}
	return nil
}
func (u *userRepo) Delete(ctx context.Context, id int64) error {
	if err := u.db.Delete(&models.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user %w", err)
	}
	return nil
}

func (u *userRepo) GetAll(ctx context.Context) ([]models.User, error) {
	users, err := gorm.G[models.User](u.db).Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users %w", err)
	}
	return users, nil
}

func (u *userRepo) GetUser(ctx context.Context, req dto.UserDetailReq) (*dto.UserDetailRes, error) {
	user, err := u.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user %w", err)
	}

	res := &dto.UserDetailRes{
		ID:           user.ID,
		Username:     user.Username,
		FullName:     user.FullName,
		Email:        user.Email,
		DepartmentID: user.DepartmentID,
		CreatedAt:    user.CreatedAt.Format("2006-01-02"),
		UpdatedAt:    user.UpdatedAt.Format("2006-01-02"),
	}
	return res, nil
}

func (u *userRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := u.db.Model(&models.User{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users %w", err)
	}
	return count, nil
}
