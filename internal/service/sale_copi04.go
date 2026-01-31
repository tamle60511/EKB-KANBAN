// service/sale_copi04_service.go

package service

import (
	"context"
	"cqs-kanban/internal/dto"
	"cqs-kanban/internal/repository"
)

type (
	saleCopi04Service struct {
		repo repository.SaleCopi04
	}

	SaleCopi04Service interface {
		GetCopi04(ctx context.Context, req dto.SaleCopi04Req) (*dto.SaleCopi04Res, error)
		GetAllCopi04(ctx context.Context) ([]dto.SaleCopi04Detail, error)
		Create(ctx context.Context, input dto.SaleCopi04Create, creator string, company string) error
		Update(ctx context.Context, id string, input dto.SaleCopi04Update, modifier string, company string) error
		Delete(ctx context.Context, id string) error
	}
)

func NewSaleCopi04Service(repo repository.SaleCopi04) SaleCopi04Service {
	return &saleCopi04Service{
		repo: repo,
	}
}

func (s *saleCopi04Service) GetCopi04(ctx context.Context, req dto.SaleCopi04Req) (*dto.SaleCopi04Res, error) {
	return s.repo.GetCopi04(ctx, req)
}

func (s *saleCopi04Service) GetAllCopi04(ctx context.Context) ([]dto.SaleCopi04Detail, error) {
	return s.repo.GetAllCopi04(ctx)
}

func (s *saleCopi04Service) Create(ctx context.Context, input dto.SaleCopi04Create, creator string, company string) error {
	return s.repo.Create(ctx, input, creator, company)
}

func (s *saleCopi04Service) Update(ctx context.Context, id string, input dto.SaleCopi04Update, modifier string, company string) error {
	return s.repo.Update(ctx, id, input, modifier, company)
}

func (s *saleCopi04Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

