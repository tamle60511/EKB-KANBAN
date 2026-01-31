package repository

import (
	"context"
	"cqs-kanban/internal/dto"

	"gorm.io/gorm"
)

type (
	copmaRepo struct {
		db *gorm.DB
	}
	CopmaRepo interface {
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

func NewCopmaRepo(db *gorm.DB) CopmaRepo {
	return &copmaRepo{
		db: db,
	}
}

func (s *copmaRepo) GetCopma(ctx context.Context) ([]dto.SaleCOPMA, error) {
	var result []dto.SaleCOPMA
	if err := s.db.WithContext(ctx).Raw(`SELECT  
            COPMA.MA001,
            COPMA.MA002,
            COPMA.MA004,
            COPMA.MA006,
            COPMA.MA015,
            COPMA.MA017,
            COPMA.MA076,
            COPMA.MA018,
            COPMA.MA019,
            COPMA.MA077,
            ISNULL(COPMA.MA002, '') AS MA002C,
            
            ISNULL(ME.ME002, '') AS MA015C,
            ISNULL(A.MR003, '') AS MA017C,
            ISNULL(C.MR003, '') AS MA018C,
            ISNULL(D.MR003, '') AS MA019C,
             ISNULL(B.MR003, '') AS MA076C,
            ISNULL(E.MR003, '') AS MA077C
        FROM COPMA
            LEFT JOIN CMSMR AS A ON A.MR002 = COPMA.MA017 AND A.MR001 = '1'
             LEFT JOIN CMSMR AS B ON A.MR002 = COPMA.MA076 AND A.MR001 = '2'
            LEFT JOIN CMSMR AS C ON C.MR002 = COPMA.MA018 AND C.MR001 = '3'
            LEFT JOIN CMSMR AS D ON D.MR002 = COPMA.MA019 AND D.MR001 = '4'
            LEFT JOIN CMSMR AS E ON E.MR002 = COPMA.MA077 AND E.MR001 = '5'
            LEFT JOIN CMSME AS ME ON ME.ME001 = COPMA.MA015`).Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
func (s *copmaRepo) GetChannel(ctx context.Context) ([]dto.Types, error) {
	var result []dto.Types
	if err := s.db.WithContext(ctx).Raw(`SELECT  
          MR002,
          MR003,
          MR004,
          MR005
        FROM CMSMR WHERE MR001 = '1'`).Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (s *copmaRepo) GetTypes(ctx context.Context) ([]dto.Types, error) {
	var result []dto.Types
	if err := s.db.WithContext(ctx).Raw(`SELECT  
          MR002,
          MR003,
          MR004,
          MR005
        FROM CMSMR WHERE MR001 = '2'`).Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (s *copmaRepo) GetRegion(ctx context.Context) ([]dto.Types, error) {
	var result []dto.Types
	if err := s.db.WithContext(ctx).Raw(`SELECT  
          MR002,
          MR003,
          MR004,
          MR005
        FROM CMSMR WHERE MR001 = '3'`).Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (s *copmaRepo) GetCountry(ctx context.Context) ([]dto.Types, error) {
	var result []dto.Types
	if err := s.db.WithContext(ctx).Raw(`SELECT  
          MR002,
          MR003,
          MR004,
          MR005
        FROM CMSMR WHERE MR001 = '4'`).Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (s *copmaRepo) GetRoute(ctx context.Context) ([]dto.Types, error) {
	var result []dto.Types
	if err := s.db.WithContext(ctx).Raw(`SELECT  
          MR002,
          MR003,
          MR004,
          MR005
        FROM CMSMR WHERE MR001 = '5'`).Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (s *copmaRepo) GetSaleDept(ctx context.Context) ([]dto.SaleDepartment, error) {
	var result []dto.SaleDepartment
	if err := s.db.WithContext(ctx).Raw(`SELECT 
  ME001,
  ME002
  FROM CMSME`).Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (s *copmaRepo) GetSaleWorkshop(ctx context.Context) ([]dto.SaleWorkshop, error) {
	var result []dto.SaleWorkshop
	if err := s.db.WithContext(ctx).Raw(`SELECT 
  MB001,
  MB002
  FROM CMSMB`).Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
func (s *copmaRepo) GetSaleItem(ctx context.Context) ([]dto.SaleItem, error) {
	var result []dto.SaleItem
	if err := s.db.WithContext(ctx).Raw(`SELECT
  MB001,
  MB002,
  MB003,
  MB004,
  MB005,
  MB006,
  MB008,
  MB017,
  MC002
  FROM INVMB
  LEFT JOIN CMSMC ON MB017=MC001`).Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
func (s *copmaRepo) GetSaleWarehouse(ctx context.Context) ([]dto.SaleWarehouse, error) {
	var result []dto.SaleWarehouse
	if err := s.db.WithContext(ctx).Raw(`SELECT
  MC001,
  MC002,
  MC003
  FROM CMSMC`).Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
func (s *copmaRepo) GetSaleMoney(ctx context.Context) ([]dto.SaleMoney, error) {
	var result []dto.SaleMoney
	if err := s.db.WithContext(ctx).Raw(`SELECT
  MF001,
  MF002
  FROM CMSMF`).Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
func (s *copmaRepo) SearchCopma(ctx context.Context, search string) ([]dto.SaleCOPMA, error) {
	var result []dto.SaleCOPMA
	query := `SELECT
		COPMA.MA001,
            COPMA.MA002,
            COPMA.MA004,
            COPMA.MA006,
            COPMA.MA015,
            COPMA.MA017,
            COPMA.MA076,
            COPMA.MA018,
            COPMA.MA019,
            COPMA.MA077,
            ISNULL(COPMA.MA002, '') AS MA002C,
            
            ISNULL(ME.ME002, '') AS MA015C,
            ISNULL(A.MR003, '') AS MA017C,
            ISNULL(C.MR003, '') AS MA018C,
            ISNULL(D.MR003, '') AS MA019C,
             ISNULL(B.MR003, '') AS MA076C,
            ISNULL(E.MR003, '') AS MA077C
        FROM COPMA
            LEFT JOIN CMSMR AS A ON A.MR002 = COPMA.MA017 AND A.MR001 = '1'
             LEFT JOIN CMSMR AS B ON A.MR002 = COPMA.MA076 AND A.MR001 = '2'
            LEFT JOIN CMSMR AS C ON C.MR002 = COPMA.MA018 AND C.MR001 = '3'
            LEFT JOIN CMSMR AS D ON D.MR002 = COPMA.MA019 AND D.MR001 = '4'
            LEFT JOIN CMSMR AS E ON E.MR002 = COPMA.MA077 AND E.MR001 = '5'
            LEFT JOIN CMSME AS ME ON ME.ME001 = COPMA.MA015 WHERE COPMA.MA001 LIKE ? OR COPMA.MA002 LIKE ?`
	likeSearch := "%" + search + "%"
	if err := s.db.WithContext(ctx).Raw(query, likeSearch, likeSearch).Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
