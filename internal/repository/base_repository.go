package repository

import (
	"context"
	"cqs-kanban/internal/dto"

	"gorm.io/gorm"
)

type (
	baseERPRepository struct {
		db *gorm.DB
	}

	BaseERPRepository interface {
		GetBaseERP(ctx context.Context, input dto.BaseERPReq) (*dto.BaseERP, error)
	}
)

func NewBaseERPRepository(ctx context.Context, db *gorm.DB) BaseERPRepository {
	return &baseERPRepository{
		db: db,
	}
}

func (r *baseERPRepository) GetBaseERP(ctx context.Context, input dto.BaseERPReq) (*dto.BaseERP, error) {
	args := map[string]any{
		"FromDate": input.FromDate,
		"ToDate":   input.ToDate,
	}

	data, err := gorm.G[map[string]any](r.db).Raw(input.SqlQuery, args).Find(ctx)
	if err != nil {
		return nil, err
	}

	convertedData := convertBytesToString(data)

	return &dto.BaseERP{
		Data: convertedData,
	}, nil
}

func convertBytesToString(data []map[string]any) []map[string]any {
	if data == nil {
		return nil
	}

	convertedData := make([]map[string]any, len(data))
	for i, row := range data {
		convertedRow := make(map[string]any)
		for key, value := range row {
			convertedRow[key] = convertValue(value)
		}
		convertedData[i] = convertedRow
	}
	return convertedData
}

func convertValue(value any) any {
	// Check if it's byte slice and convert to string
	if bytes, ok := value.([]byte); ok {
		return string(bytes)
	}
	return value
}
