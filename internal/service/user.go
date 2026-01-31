package service

import (
	"context"
	"cqs-kanban/internal/dto"
	"cqs-kanban/internal/repository"
	"cqs-kanban/internal/utils"
	"errors"
	"fmt"
)

type (
	userService struct {
		userRepo repository.UserRepo
	}
	UserService interface {
		Create(ctx context.Context, req dto.UserCreateReq) error
		GetByID(ctx context.Context, req dto.UserDetailReq) (*dto.UserDetailRes, error)
		GetByUsername(ctx context.Context, username string) (*dto.UserDetailRes, error)
		Update(ctx context.Context, id int, req dto.UserUpdateReq) error
		UpdatePassword(ctx context.Context, id int64, req dto.UpdatePasswordRequest) error
		Delete(ctx context.Context, id int) error
		GetAll(ctx context.Context) ([]dto.UserDetailRes, error)
		Count(ctx context.Context) (int64, error)
	}
)

func NewUserService(userRepo repository.UserRepo) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (u *userService) Create(ctx context.Context, req dto.UserCreateReq) error {
	return u.userRepo.Create(ctx, req)
}
func (u *userService) GetByID(ctx context.Context, req dto.UserDetailReq) (*dto.UserDetailRes, error) {
	return u.userRepo.GetUser(ctx, req)
}
func (u *userService) GetByUsername(ctx context.Context, username string) (*dto.UserDetailRes, error) {
	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	return &dto.UserDetailRes{
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
	}, nil
}
func (u *userService) Update(ctx context.Context, id int, req dto.UserUpdateReq) error {
	return u.userRepo.Update(ctx, int64(id), req)
}
func (u *userService) UpdatePassword(ctx context.Context, id int64, req dto.UpdatePasswordRequest) error {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}
	if !utils.CheckPasswordHash(req.CurrentPassword, user.Password) {
		return errors.New("current password is incorrect")
	}
	if req.NewPassword != req.ConfirmPassword {
		return errors.New("new password cannot be the same as the current password")
	}
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}
	return u.userRepo.UpdatePassword(ctx, id, hashedPassword)
}
func (u *userService) Delete(ctx context.Context, id int) error {
	return u.userRepo.Delete(ctx, int64(id))
}
func (u *userService) GetAll(ctx context.Context) ([]dto.UserDetailRes, error) {
	users, err := u.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var userDetails []dto.UserDetailRes
	for _, user := range users {
		userDetails = append(userDetails, dto.UserDetailRes{
			ID:           user.ID,
			Username:     user.Username,
			FullName:     user.FullName,
			Email:        user.Email,
			DepartmentID: user.DepartmentID,
			CreatedAt:    user.CreatedAt.Format("2006-01-02"),
			UpdatedAt:    user.UpdatedAt.Format("2006-01-02"),
		})
	}

	return userDetails, nil
}

func (u *userService) Count(ctx context.Context) (int64, error) {
	count, err := u.userRepo.Count(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}
