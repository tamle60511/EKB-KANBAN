package service

import (
	"context"
	"cqs-kanban/internal/dto"
	"cqs-kanban/internal/repository"
)

type (
	copmaService struct {
		repo repository.CopmaRepo
	}
	CopmaService interface {
		GetCopma(ctx context.Context) ([]dto.SaleCOPMA, error)
		GetChannel(ctx context.Context) ([]dto.Types, error)
		GetTypes(ctx context.Context) ([]dto.Types, error)
		GetRegion(ctx context.Context) ([]dto.Types, error)
		GetCountry(ctx context.Context) ([]dto.Types, error)
		GetRoute(ctx context.Context) ([]dto.Types, error)
		GetSaleDept(ctx context.Context) ([]dto.SaleDepartment, error)
		GetSaleWorkshop(ctx context.Context) ([]dto.SaleWorkshop, error)
		GetSaleItem(ctx context.Context) ([]dto.SaleItem, error)
		GetSaleWarehouse(ctx context.Context) ([]dto.SaleWarehouse, error)
		GetSaleMoney(ctx context.Context) ([]dto.SaleMoney, error)
		SearchCopma(ctx context.Context, search string) ([]dto.SaleCOPMA, error)
	}
)

func NewCopmaService(repo repository.CopmaRepo) CopmaService {
	return &copmaService{
		repo: repo,
	}
}

func (s *copmaService) GetCopma(ctx context.Context) ([]dto.SaleCOPMA, error) {
	return s.repo.GetCopma(ctx)
}
func (s *copmaService) GetChannel(ctx context.Context) ([]dto.Types, error) {
	return s.repo.GetChannel(ctx)
}
func (s *copmaService) GetTypes(ctx context.Context) ([]dto.Types, error) {
	return s.repo.GetTypes(ctx)
}
func (s *copmaService) GetRegion(ctx context.Context) ([]dto.Types, error) {
	return s.repo.GetRegion(ctx)
}
func (s *copmaService) GetCountry(ctx context.Context) ([]dto.Types, error) {
	return s.repo.GetCountry(ctx)
}
func (s *copmaService) GetRoute(ctx context.Context) ([]dto.Types, error) {
	return s.repo.GetRoute(ctx)
}

func (s *copmaService) GetSaleDept(ctx context.Context) ([]dto.SaleDepartment, error) {
	return s.repo.GetSaleDept(ctx)
}

func (s *copmaService) GetSaleWorkshop(ctx context.Context) ([]dto.SaleWorkshop, error) {
	return s.repo.GetSaleWorkshop(ctx)
}
func (s *copmaService) GetSaleItem(ctx context.Context) ([]dto.SaleItem, error) {
	return s.repo.GetSaleItem(ctx)
}

func (s *copmaService) GetSaleWarehouse(ctx context.Context) ([]dto.SaleWarehouse, error) {
	return s.repo.GetSaleWarehouse(ctx)
}

func (s *copmaService) GetSaleMoney(ctx context.Context) ([]dto.SaleMoney, error) {
	return s.repo.GetSaleMoney(ctx)
}

func (s *copmaService) SearchCopma(ctx context.Context, search string) ([]dto.SaleCOPMA, error) {
	return s.repo.SearchCopma(ctx, search)
}
