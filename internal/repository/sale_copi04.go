// repository/sale_copi04_repository.go

package repository

import (
	"context"
	"cqs-kanban/internal/dto"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	saleCopi04 struct {
		db *gorm.DB
	}

	SaleCopi04 interface {
		GetCopi04(ctx context.Context, req dto.SaleCopi04Req) (*dto.SaleCopi04Res, error)
		GetAllCopi04(ctx context.Context) ([]dto.SaleCopi04Detail, error)
		Create(ctx context.Context, input dto.SaleCopi04Create, creator string, company string) error
		Update(ctx context.Context, id string, input dto.SaleCopi04Update, modifier string, company string) error
		Delete(ctx context.Context, id string) error
	}
)

func NewSaleCopi04(db *gorm.DB) SaleCopi04 {
	return &saleCopi04{
		db: db,
	}
}

func (s *saleCopi04) GetCopi04(ctx context.Context, req dto.SaleCopi04Req) (*dto.SaleCopi04Res, error) {
	var header dto.SaleCopi04Model

	headerQuery := `
		SELECT 
			RTRIM(COPME.ME001) AS ME001,
			COPME.ME002 AS ME002,    
			COPME.ME003 AS ME003,    
			COPME.ME004 AS ME004,    
			COPME.ME005 AS ME005,    
			COPME.ME006 AS ME006,    
			RTRIM(COPME.ME007) AS ME007,   
			RTRIM(COPME.ME008) AS ME008,    
			RTRIM(COPME.ME009) AS ME009,    
			RTRIM(COPME.ME010) AS ME010,   
			RTRIM(COPME.ME011) AS ME011,    
			RTRIM(COPME.ME012) AS ME012,   
			RTRIM(COPME.ME013) AS ME013,   
			RTRIM(COPME.ME014) AS ME014
		FROM COPME AS COPME
		WHERE RTRIM(COPME.ME001) = ?
	`

	if err := s.db.WithContext(ctx).Raw(headerQuery, strings.TrimSpace(req.ME001)).Scan(&header).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("record not found")
		}
		return nil, err
	}

	var details []dto.SeleCopi04ColModel

	detailQuery := `
		SELECT 
		RTRIM(COPMF.MF001) AS MF001,
		RTRIM(COPMF.MF002) AS MF002,
		RTRIM(COPMF.MF003) AS MF003,    
		RTRIM(COPMF.MF004) AS MF004,    
		RTRIM(COPMF.MF005) AS MF005,    
		CASE 
			WHEN ISDATE(COPMF.MF006) = 1 THEN CONVERT(VARCHAR(10), CONVERT(DATETIME, COPMF.MF006), 23) 
			ELSE NULL 
		END AS MF006,    
		RTRIM(COPMF.MF007) AS MF007,   
		COPMF.MF008 AS MF008,    
		COPMF.MF009 AS MF009,    
		RTRIM(COPMF.MF010) AS MF010,   
		RTRIM(COPMF.MF011) AS MF011,    
		COPMF.MF012 AS MF012,   
		RTRIM(COPMF.MF013) AS MF013,   
		COPMF.MF014 AS MF014,   
		COPMF.MF015 AS MF015,
		RTRIM(COPMF.MF020) AS MF020
	FROM COPMF AS COPMF
	WHERE RTRIM(COPMF.MF001) = ?
	`

	if err := s.db.WithContext(ctx).Raw(detailQuery, strings.TrimSpace(req.ME001)).Scan(&details).Error; err != nil {
		return nil, err
	}

	return &dto.SaleCopi04Res{
		Header: header,
		Detail: details,
	}, nil
}

func (s *saleCopi04) GetAllCopi04(ctx context.Context) ([]dto.SaleCopi04Detail, error) {
	var results []dto.SaleCopi04Detail

	query := `
		SELECT 
			COPME.ME001 AS ME001,     
			COPME.ME002 AS ME002,     
			COPME.ME003 AS ME003,     
			COPME.ME004 AS ME004,     
			COPME.ME005 AS ME005,     
			COPME.ME006 AS ME006, 
			COPME.ME008 AS ME008,     
			COPME.ME009 AS ME009,     
			COPME.ME010 AS ME010,     
			COPME.ME011 AS ME011,    
			COPME.ME013 AS ME012,    
			ISNULL(MA002, '') AS ME002C,
			ISNULL(A.MR003, '') AS ME003C,
			ISNULL(C.MR003, '') AS ME004C,
			ISNULL(D.MR003, '') AS ME005C,
			ISNULL(CMSME.ME002, '') AS ME006C,
			ISNULL(B.MR003, '') AS ME010C,
			ISNULL(E.MR003, '') AS ME011C,
			ISNULL(MB002, '') AS ME012C
		FROM COPME AS COPME
		LEFT JOIN COPMA AS COPMA ON COPMA.MA001 = COPME.ME002
		LEFT JOIN CMSMR AS A ON A.MR002 = COPME.ME003 AND A.MR001 = '1'
		LEFT JOIN CMSMR AS B ON B.MR002 = COPME.ME010 AND B.MR001 = '2'
		LEFT JOIN CMSMR AS C ON C.MR002 = COPME.ME004 AND C.MR001 = '3'
		LEFT JOIN CMSMR AS D ON D.MR002 = COPME.ME005 AND D.MR001 = '4'
		LEFT JOIN CMSMR AS E ON E.MR002 = COPME.ME011 AND E.MR001 = '5'
		LEFT JOIN CMSME AS CMSME ON CMSME.ME001 = COPME.ME006
		LEFT JOIN CMSMB AS CMSMB ON CMSMB.MB001 = COPME.ME013 
		WHERE COPME.ME014 = '1' 
	`

	if err := s.db.WithContext(ctx).Raw(query).Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (s *saleCopi04) Create(ctx context.Context, input dto.SaleCopi04Create, creator string, company string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if input.ME001 == "" {
			return fmt.Errorf("ME001 cannot be empty")
		}

		var count int64
		if err := tx.Model(&dto.SaleCopi04Model{}).
			Where("ME001 = ?", input.ME001).
			Count(&count).Error; err != nil {
			return fmt.Errorf("failed to check duplicate: %w", err)
		}

		if count > 0 {
			return fmt.Errorf("ME001 '%s' already exists", input.ME001)
		}

		header := dto.SaleCopi04Model{
			ME001:       input.ME001,
			ME002:       input.ME002,
			ME003:       input.ME003,
			ME004:       input.ME004,
			ME005:       input.ME005,
			ME006:       input.ME006,
			ME007:       "", // Check if can be left empty
			ME008:       input.ME008,
			ME009:       "", // Check if can be left empty
			ME010:       "", // Check if can be left empty
			ME011:       input.ME011,
			ME012:       "", // Check if can be left empty
			ME013:       input.ME013,
			ME014:       "1", // Active record
			COMPANY:     company,
			CREATOR:     creator,
			USR_GROUP:   "", // Check if can be left empty
			CREATE_DATE: time.Now().Format("20060102"),
			CREATE_TIME: time.Now().Format("15:04:05"), // Add a proper time value
			MODIFIER:    "",
			MODI_DATE:   "",
			MODI_TIME:   "",
			CREATE_AP:   "",
			CREATE_PRID: "",
			MODI_AP:     "",
			MODI_PRID:   "",
		}

		log.Printf("Inserting COPME: %+v", header)

		if err := tx.Create(&header).Error; err != nil {
			return fmt.Errorf("failed to create header: %w", err)
		}

		// Check if Columns is provided, process only if not empty
		if len(input.Columns) > 0 {
			for i, v := range input.Columns {
				detail := dto.SeleCopi04Col{
					MF001:       input.ME001,
					MF002:       v.MF002,
					MF003:       v.MF003,
					MF004:       v.MF004,
					MF005:       v.MF005,
					MF006:       v.MF006,
					MF007:       v.MF007,
					MF008:       v.MF008,
					MF009:       v.MF009,
					MF010:       v.MF010,
					MF011:       v.MF011,
					MF012:       v.MF012,
					MF013:       v.MF013,
					MF014:       v.MF014,
					MF015:       v.MF015,
					CREATOR:     creator,
					COMPANY:     company,
					CREATE_DATE: time.Now().Format("20060102"),
				}

				if err := tx.Create(&detail).Error; err != nil {
					return fmt.Errorf("failed to create detail row %d: %w", i+1, err)
				}
			}
		} else {
			log.Printf("No columns provided, skipping detail insertion")
		}

		return nil
	})
}
func (s *saleCopi04) Update(ctx context.Context, id string, input dto.SaleCopi04Update, modifier string, company string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Cập nhật header
		headerUpdates := map[string]interface{}{
			"ME002": input.ME002,
			"ME003": input.ME003,
			"ME004": input.ME004,
			"ME005": input.ME005,
			"ME006": input.ME006,
			"ME007": input.ME007,
			"ME008": input.ME008,
			"ME009": input.ME009,
			"ME010": input.ME010,

			"ME011":     input.ME011,
			"ME012":     input.ME012,
			"ME013":     input.ME013,
			"MODIFIER":  modifier,
			"MODI_DATE": time.Now().Format("20060102"),
			"MODI_TIME": time.Now().Format("15:04:05"),
		}

		if err := tx.Model(&dto.SaleCopi04Model{}).
			Where("ME001 = ?", id).
			Updates(headerUpdates).Error; err != nil {
			return fmt.Errorf("failed to update header: %w", err)
		}

		// Xóa các hàng chi tiết cũ
		if err := tx.Where("MF001 = ?", id).Delete(&dto.SeleCopi04Col{}).Error; err != nil {
			return fmt.Errorf("failed to delete existing detail rows: %w", err)
		}

		// Tạo các hàng chi tiết mới
		for i, v := range input.Columns {
			detail := dto.SeleCopi04Col{
				MF001:       id, // Yêu cầu ở đây phải là id của header
				MF002:       v.MF002,
				MF003:       v.MF003,
				MF004:       v.MF004,
				MF005:       v.MF005,
				MF006:       v.MF006,
				MF007:       v.MF007,
				MF008:       v.MF008,
				MF009:       v.MF009,
				MF010:       v.MF010,
				MF011:       v.MF011,
				MF012:       v.MF012,
				MF013:       v.MF013,
				MF014:       v.MF014,
				MF015:       v.MF015,
				CREATOR:     modifier,
				COMPANY:     company,
				CREATE_DATE: time.Now().Format("20060102"),
			}

			if err := tx.Create(&detail).Error; err != nil {
				return fmt.Errorf("failed to create detail row %d: %w", i+1, err)
			}
		}

		return nil
	})
}

func (s *saleCopi04) Delete(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete columns first (foreign key)
		if err := tx.Where("MF001 = ?", id).Delete(&dto.SeleCopi04Col{}).Error; err != nil {
			return fmt.Errorf("failed to delete columns: %w", err)
		}

		// Delete header
		if err := tx.Where("ME001 = ?", id).Delete(&dto.SaleCopi04Model{}).Error; err != nil {
			return fmt.Errorf("failed to delete header: %w", err)
		}

		return nil
	})
}

