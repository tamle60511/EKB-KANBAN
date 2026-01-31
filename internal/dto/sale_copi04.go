package dto

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"
)

type SaleCopi04Req struct {
	ME001 string `json:"ME001" uri:"id" validate:"required"`
}

type SaleCopi04Res struct {
	Header SaleCopi04Model      `json:"header"`
	Detail []SeleCopi04ColModel `json:"detail"`
}

type SaleCopi04Create struct {
	ME001   string           `json:"ME001" validate:"required,nchar=100"`
	ME002   string           `json:"ME002,omitempty" validate:"nchar=10"`
	ME003   string           `json:"ME003,omitempty" validate:"nchar=6"`
	ME004   string           `json:"ME004,omitempty" validate:"nchar=6"`
	ME005   string           `json:"ME005,omitempty" validate:"nchar=6"`
	ME006   string           `json:"ME006,omitempty" validate:"nchar=10"`
	ME007   string           `json:"ME007,omitempty" validate:"nchar=10"`
	ME008   string           `json:"ME008,omitempty" validate:"nchar=1"`
	ME009   string           `json:"ME009,omitempty" validate:"nchar=255"`
	ME010   string           `json:"ME010,omitempty" validate:"nchar=6"`
	ME011   string           `json:"ME011,omitempty" validate:"nchar=6"`
	ME012   string           `json:"ME012,omitempty" validate:"nchar=6"`
	ME013   string           `json:"ME013,omitempty" validate:"nchar=6"`
	ME014   string           `json:"ME014,omitempty" validate:"nchar=1"`
	Columns []*SeleCopi04Col `json:"columns" validate:"dive"`
}

// Validate tùy chỉnh nếu không dùng custom validator
func (c *SaleCopi04Create) ValidateLength() error {
	if utf8.RuneCountInString(c.ME001) > 100 {
		return fmt.Errorf("ME001 exceeds maximum length of 100 characters")
	}
	// Tương tự cho các field khác...
	return nil
}

type SaleCopi04Update struct {
	ME002   *string          `json:"ME002" validate:"omitempty,max=10"`
	ME003   *string          `json:"ME003" validate:"omitempty,max=6"`
	ME004   *string          `json:"ME004" validate:"omitempty,max=6"`
	ME005   *string          `json:"ME005" validate:"omitempty,max=6"`
	ME006   *string          `json:"ME006" validate:"omitempty,max=10"`
	ME007   *string          `json:"ME007" validate:"omitempty,max=10"`
	ME008   *string          `json:"ME008" validate:"omitempty,max=1"`
	ME009   *string          `json:"ME009" validate:"omitempty,max=255"`
	ME010   *string          `json:"ME010" validate:"omitempty,max=6"`
	ME011   *string          `json:"ME011" validate:"omitempty,max=6"`
	ME012   *string          `json:"ME012" validate:"omitempty,max=6"`
	ME013   *string          `json:"ME013" validate:"omitempty,max=6"`
	Columns []*SeleCopi04Col `json:"columns"`
}

type SaleCopi04Model struct {
	ME001       string `json:"ME001" gorm:"column:ME001;primaryKey;type:nchar(100)"`
	ME002       string `json:"ME002,omitempty" gorm:"column:ME002;type:nvarchar(10)"`
	ME003       string `json:"ME003,omitempty" gorm:"column:ME003;type:nvarchar(6)"`
	ME004       string `json:"ME004,omitempty" gorm:"column:ME004;type:nvarchar(6)"`
	ME005       string `json:"ME005,omitempty" gorm:"column:ME005;type:nvarchar(6)"`
	ME006       string `json:"ME006,omitempty" gorm:"column:ME006;type:nvarchar(10)"`
	ME007       string `json:"ME007,omitempty" gorm:"column:ME007;type:nvarchar(10)"`
	ME008       string `json:"ME008,omitempty" gorm:"column:ME008;type:nvarchar(1)"`
	ME009       string `json:"ME009,omitempty" gorm:"column:ME009;type:nvarchar(255)"`
	ME010       string `json:"ME010,omitempty" gorm:"column:ME010;type:nvarchar(6)"`
	ME011       string `json:"ME011,omitempty" gorm:"column:ME011;type:nvarchar(6)"`
	ME012       string `json:"ME012,omitempty" gorm:"column:ME012;type:nvarchar(6)"`
	ME013       string `json:"ME013,omitempty" gorm:"column:ME013;type:nvarchar(6)"`
	ME014       string `json:"ME014,omitempty" gorm:"column:ME014;type:nvarchar(1)"`
	COMPANY     string `json:"COMPANY" gorm:"column:COMPANY;type:nvarchar(20)"` // Adjusted
	CREATOR     string `json:"CREATOR" gorm:"column:CREATOR;type:nvarchar(10)"` // Adjusted
	USR_GROUP   string `json:"USR_GROUP,omitempty" gorm:"column:USR_GROUP;type:nvarchar(10)"`
	CREATE_DATE string `json:"CREATE_DATE,omitempty" gorm:"column:CREATE_DATE;type:nvarchar(8)"`  // Adjusted
	CREATE_TIME string `json:"CREATE_TIME,omitempty" gorm:"column:CREATE_TIME;type:nvarchar(20)"` // Adjusted
	MODIFIER    string `json:"MODIFIER" gorm:"column:MODIFIER;type:nvarchar(10)"`                 // Adjusted
	MODI_DATE   string `json:"MODI_DATE,omitempty" gorm:"column:MODI_DATE;type:nvarchar(8)"`      // Adjusted
	MODI_TIME   string `json:"MODI_TIME,omitempty" gorm:"column:MODI_TIME;type:nvarchar(20)"`     // Adjusted
	CREATE_AP   string `json:"CREATE_AP,omitempty" gorm:"column:CREATE_AP;type:nvarchar(255)"`    // New field
	CREATE_PRID string `json:"CREATE_PRID,omitempty" gorm:"column:CREATE_PRID;type:nvarchar(50)"` // New field
	MODI_AP     string `json:"MODI_AP,omitempty" gorm:"column:MODI_AP;type:nvarchar(255)"`        // New field
	MODI_PRID   string `json:"MODI_PRID,omitempty" gorm:"column:MODI_PRID;type:nvarchar(50)"`
}

func (SaleCopi04Model) TableName() string {
	return "COPME"
}

func (m SaleCopi04Model) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ME001 string `json:"ME001"`
		ME002 string `json:"ME002,omitempty"`
		ME003 string `json:"ME003,omitempty"`
		ME004 string `json:"ME004,omitempty"`
		ME005 string `json:"ME005,omitempty"`
		ME006 string `json:"ME006,omitempty"`
		ME007 string `json:"ME007,omitempty"`
		ME008 string `json:"ME008,omitempty"`
		ME009 string `json:"ME009,omitempty"`
		ME010 string `json:"ME010,omitempty"`
		ME011 string `json:"ME011,omitempty"`
		ME012 string `json:"ME012,omitempty"`
		ME013 string `json:"ME013,omitempty"`
		ME014 string `json:"ME014,omitempty"`
	}{
		ME001: strings.TrimSpace(m.ME001),
		ME002: strings.TrimSpace(m.ME002),
		ME003: strings.TrimSpace(m.ME003),
		ME004: strings.TrimSpace(m.ME004),
		ME005: strings.TrimSpace(m.ME005),
		ME006: strings.TrimSpace(m.ME006),
		ME007: strings.TrimSpace(m.ME007),
		ME008: strings.TrimSpace(m.ME008),
		ME009: strings.TrimSpace(m.ME009),
		ME010: strings.TrimSpace(m.ME010),
		ME011: strings.TrimSpace(m.ME011),
		ME012: strings.TrimSpace(m.ME012),
		ME013: strings.TrimSpace(m.ME013),
		ME014: strings.TrimSpace(m.ME014),
	})
}

type SaleCopi04Detail struct {
	ME001  string `json:"ME001"`
	ME002  string `json:"ME002"`
	ME003  string `json:"ME003"`
	ME004  string `json:"ME004"`
	ME005  string `json:"ME005"`
	ME006  string `json:"ME006"`
	ME007  string `json:"ME007"`
	ME008  string `json:"ME008"`
	ME009  string `json:"ME009"`
	ME010  string `json:"ME010"`
	ME011  string `json:"ME011"`
	ME012  string `json:"ME012"`
	ME013  string `json:"ME013"`
	ME002C string `json:"ME002C"`
	ME003C string `json:"ME003C"`
	ME004C string `json:"ME004C"`
	ME005C string `json:"ME005C"`
	ME006C string `json:"ME006C"`
	ME010C string `json:"ME010C"`
	ME011C string `json:"ME011C"`
	ME012C string `json:"ME012C"`
}

func (m SaleCopi04Detail) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ME001  string `json:"ME001"`
		ME002  string `json:"ME002"`
		ME003  string `json:"ME003"`
		ME004  string `json:"ME004"`
		ME005  string `json:"ME005"`
		ME006  string `json:"ME006"`
		ME007  string `json:"ME007"`
		ME008  string `json:"ME008"`
		ME009  string `json:"ME009"`
		ME010  string `json:"ME010"`
		ME011  string `json:"ME011"`
		ME012  string `json:"ME012"`
		ME013  string `json:"ME013"`
		ME002C string `json:"ME002C"`
		ME003C string `json:"ME003C"`
		ME004C string `json:"ME004C"`
		ME005C string `json:"ME005C"`
		ME006C string `json:"ME006C"`
		ME010C string `json:"ME010C"`
		ME011C string `json:"ME011C"`
		ME012C string `json:"ME012C"`
	}{
		ME001:  strings.TrimSpace(m.ME001),
		ME002:  strings.TrimSpace(m.ME002),
		ME003:  strings.TrimSpace(m.ME003),
		ME004:  strings.TrimSpace(m.ME004),
		ME005:  strings.TrimSpace(m.ME005),
		ME006:  strings.TrimSpace(m.ME006),
		ME007:  strings.TrimSpace(m.ME007),
		ME008:  strings.TrimSpace(m.ME008),
		ME009:  strings.TrimSpace(m.ME009),
		ME010:  strings.TrimSpace(m.ME010),
		ME011:  strings.TrimSpace(m.ME011),
		ME012:  strings.TrimSpace(m.ME012),
		ME013:  strings.TrimSpace(m.ME013),
		ME002C: strings.TrimSpace(m.ME002C),
		ME003C: strings.TrimSpace(m.ME003C),
		ME004C: strings.TrimSpace(m.ME004C),
		ME005C: strings.TrimSpace(m.ME005C),
		ME006C: strings.TrimSpace(m.ME006C),
		ME010C: strings.TrimSpace(m.ME010C),
		ME011C: strings.TrimSpace(m.ME011C),
		ME012C: strings.TrimSpace(m.ME012C),
	})
}

type SeleCopi04Col struct {
	MF001       string  `json:"MF001" gorm:"column:MF001;primaryKey"`
	MF002       string  `json:"MF002" gorm:"column:MF002;primaryKey"`
	MF003       string  `json:"MF003,omitempty" gorm:"column:MF003"`
	MF004       string  `json:"MF004,omitempty" gorm:"column:MF004"`
	MF005       string  `json:"MF005,omitempty" gorm:"column:MF005"`
	MF006       string  `json:"MF006,omitempty" gorm:"column:MF006"`
	MF007       string  `json:"MF007,omitempty" gorm:"column:MF007"`
	MF008       float64 `json:"MF008,omitempty" gorm:"column:MF008"`
	MF009       float64 `json:"MF009,omitempty" gorm:"column:MF009"`
	MF010       string  `json:"MF010,omitempty" gorm:"column:MF010"`
	MF011       string  `json:"MF011,omitempty" gorm:"column:MF011"`
	MF012       float64 `json:"MF012,omitempty" gorm:"column:MF012"`
	MF013       string  `json:"MF013,omitempty" gorm:"column:MF013"`
	MF014       float64 `json:"MF014,omitempty" gorm:"column:MF014"`
	MF015       int     `json:"MF015,omitempty" gorm:"column:MF015"`
	MF016       string  `json:"MF016,omitempty" gorm:"column:MF016"`
	MF017       string  `json:"MF017,omitempty" gorm:"column:MF017"`
	MF018       string  `json:"MF018,omitempty" gorm:"column:MF018"`
	MF019       string  `json:"MF019,omitempty" gorm:"column:MF019"`
	MF020       string  `json:"MF020,omitempty" gorm:"column:MF020"`
	MF021       float64 `json:"MF021,omitempty" gorm:"column:MF021"`
	MF022       float64 `json:"MF022,omitempty" gorm:"column:MF022"`
	MF023       string  `json:"MF023,omitempty" gorm:"column:MF023"`
	MF024       string  `json:"MF024,omitempty" gorm:"column:MF024"`
	MF025       string  `json:"MF025,omitempty" gorm:"column:MF025"`
	MF026       string  `json:"MF026,omitempty" gorm:"column:MF026"`
	MF027       string  `json:"MF027,omitempty" gorm:"column:MF027"`
	MF028       string  `json:"MF028,omitempty" gorm:"column:MF028"`
	MF029       string  `json:"MF029,omitempty" gorm:"column:MF029"`
	MF030       string  `json:"MF030,omitempty" gorm:"column:MF030"`
	MF031       string  `json:"MF031,omitempty" gorm:"column:MF031"`
	MF032       string  `json:"MF032,omitempty" gorm:"column:MF032"`
	MF033       string  `json:"MF033,omitempty" gorm:"column:MF033"`
	UDF01       string  `json:"UDF01,omitempty" gorm:"column:UDF01"`
	UDF02       string  `json:"UDF02,omitempty" gorm:"column:UDF02"`
	UDF03       string  `json:"UDF03,omitempty" gorm:"column:UDF03"`
	UDF04       string  `json:"UDF04,omitempty" gorm:"column:UDF04"`
	UDF05       string  `json:"UDF05,omitempty" gorm:"column:UDF05"`
	UDF06       float64 `json:"UDF06,omitempty" gorm:"column:UDF06"`
	UDF07       float64 `json:"UDF07,omitempty" gorm:"column:UDF07"`
	UDF08       float64 `json:"UDF08,omitempty" gorm:"column:UDF08"`
	UDF09       float64 `json:"UDF09,omitempty" gorm:"column:UDF09"`
	UDF10       float64 `json:"UDF10,omitempty" gorm:"column:UDF10"`
	COMPANY     string  `json:"COMPANY" gorm:"column:COMPANY;type:nvarchar(20)"` // Adjusted
	CREATOR     string  `json:"CREATOR" gorm:"column:CREATOR;type:nvarchar(10)"` // Adjusted
	USR_GROUP   string  `json:"USR_GROUP,omitempty" gorm:"column:USR_GROUP;type:nvarchar(10)"`
	CREATE_DATE string  `json:"CREATE_DATE,omitempty" gorm:"column:CREATE_DATE;type:nvarchar(8)"`  // Adjusted
	CREATE_TIME string  `json:"CREATE_TIME,omitempty" gorm:"column:CREATE_TIME;type:nvarchar(20)"` // Adjusted
	MODIFIER    string  `json:"MODIFIER" gorm:"column:MODIFIER;type:nvarchar(10)"`                 // Adjusted
	MODI_DATE   string  `json:"MODI_DATE,omitempty" gorm:"column:MODI_DATE;type:nvarchar(8)"`      // Adjusted
	MODI_TIME   string  `json:"MODI_TIME,omitempty" gorm:"column:MODI_TIME;type:nvarchar(20)"`     // Adjusted
	CREATE_AP   string  `json:"CREATE_AP,omitempty" gorm:"column:CREATE_AP;type:nvarchar(255)"`    // New field
	CREATE_PRID string  `json:"CREATE_PRID,omitempty" gorm:"column:CREATE_PRID;type:nvarchar(50)"` // New field
	MODI_AP     string  `json:"MODI_AP,omitempty" gorm:"column:MODI_AP;type:nvarchar(255)"`        // New field
	MODI_PRID   string  `json:"MODI_PRID,omitempty" gorm:"column:MODI_PRID;type:nvarchar(50)"`
}

func (SeleCopi04Col) TableName() string {
	return "COPMF"
}

type SeleCopi04ColModel struct {
	MF001 string  `json:"MF001" gorm:"column:MF001;primaryKey"`
	MF002 string  `json:"MF002" gorm:"column:MF002;primaryKey"`
	MF003 string  `json:"MF003,omitempty" gorm:"column:MF003"`
	MF004 string  `json:"MF004,omitempty" gorm:"column:MF004"`
	MF005 string  `json:"MF005,omitempty" gorm:"column:MF005"`
	MF006 string  `json:"MF006,omitempty" gorm:"column:MF006"`
	MF007 string  `json:"MF007,omitempty" gorm:"column:MF007"`
	MF008 float64 `json:"MF008,omitempty" gorm:"column:MF008"`
	MF009 float64 `json:"MF009,omitempty" gorm:"column:MF009"`
	MF010 string  `json:"MF010,omitempty" gorm:"column:MF010"`
	MF011 string  `json:"MF011,omitempty" gorm:"column:MF011"`
	MF012 float64 `json:"MF012,omitempty" gorm:"column:MF012"`
	MF013 string  `json:"MF013,omitempty" gorm:"column:MF013"`
	MF014 float64 `json:"MF014,omitempty" gorm:"column:MF014"`
	MF015 int     `json:"MF015,omitempty" gorm:"column:MF015"`
	MF016 string  `json:"MF016,omitempty" gorm:"column:MF016"`
	MF017 string  `json:"MF017,omitempty" gorm:"column:MF017"`
	MF018 string  `json:"MF018,omitempty" gorm:"column:MF018"`
	MF019 string  `json:"MF019,omitempty" gorm:"column:MF019"`
	MF020 string  `json:"MF020,omitempty" gorm:"column:MF020"`
	MF021 float64 `json:"MF021,omitempty" gorm:"column:MF021"`
	MF022 float64 `json:"MF022,omitempty" gorm:"column:MF022"`
	MF023 string  `json:"MF023,omitempty" gorm:"column:MF023"`
	MF024 string  `json:"MF024,omitempty" gorm:"column:MF024"`
	MF025 string  `json:"MF025,omitempty" gorm:"column:MF025"`
	MF026 string  `json:"MF026,omitempty" gorm:"column:MF026"`
	MF027 string  `json:"MF027,omitempty" gorm:"column:MF027"`
	MF028 string  `json:"MF028,omitempty" gorm:"column:MF028"`
	MF029 string  `json:"MF029,omitempty" gorm:"column:MF029"`
	MF030 string  `json:"MF030,omitempty" gorm:"column:MF030"`
	MF031 string  `json:"MF031,omitempty" gorm:"column:MF031"`
	MF032 string  `json:"MF032,omitempty" gorm:"column:MF032"`
	MF033 string  `json:"MF033,omitempty" gorm:"column:MF033"`
	UDF01 string  `json:"UDF01,omitempty" gorm:"column:UDF01"`
	UDF02 string  `json:"UDF02,omitempty" gorm:"column:UDF02"`
	UDF03 string  `json:"UDF03,omitempty" gorm:"column:UDF03"`
	UDF04 string  `json:"UDF04,omitempty" gorm:"column:UDF04"`
	UDF05 string  `json:"UDF05,omitempty" gorm:"column:UDF05"`
	UDF06 float64 `json:"UDF06,omitempty" gorm:"column:UDF06"`
	UDF07 float64 `json:"UDF07,omitempty" gorm:"column:UDF07"`
	UDF08 float64 `json:"UDF08,omitempty" gorm:"column:UDF08"`
	UDF09 float64 `json:"UDF09,omitempty" gorm:"column:UDF09"`
	UDF10 float64 `json:"UDF10,omitempty" gorm:"column:UDF10"`
}

func (m SeleCopi04ColModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		MF001 string  `json:"MF001"`
		MF002 string  `json:"MF002"`
		MF003 string  `json:"MF003,omitempty"`
		MF004 string  `json:"MF004,omitempty"`
		MF005 string  `json:"MF005,omitempty"`
		MF006 string  `json:"MF006,omitempty"`
		MF007 string  `json:"MF007,omitempty"`
		MF008 float64 `json:"MF008,omitempty"`
		MF009 float64 `json:"MF009,omitempty"`
		MF010 string  `json:"MF010,omitempty"`
		MF011 string  `json:"MF011,omitempty"`
		MF012 float64 `json:"MF012,omitempty"`
		MF013 string  `json:"MF013,omitempty"`
		MF014 float64 `json:"MF014,omitempty"`
		MF015 int     `json:"MF015,omitempty"`
		MF016 string  `json:"MF016,omitempty"`
		MF017 string  `json:"MF017,omitempty"`
		MF018 string  `json:"MF018,omitempty"`
		MF019 string  `json:"MF019,omitempty"`
		MF020 string  `json:"MF020,omitempty"`
		MF021 float64 `json:"MF021,omitempty"`
		MF022 float64 `json:"MF022,omitempty"`
		MF023 string  `json:"MF023,omitempty"`
		MF024 string  `json:"MF024,omitempty"`
		MF025 string  `json:"MF025,omitempty"`
		MF026 string  `json:"MF026,omitempty"`
		MF027 string  `json:"MF027,omitempty"`
		MF028 string  `json:"MF028,omitempty"`
		MF029 string  `json:"MF029,omitempty"`
		MF030 string  `json:"MF030,omitempty"`
		MF031 string  `json:"MF031,omitempty"`
		MF032 string  `json:"MF032,omitempty"`
		MF033 string  `json:"MF033,omitempty"`
		UDF01 string  `json:"UDF01,omitempty"`
		UDF02 string  `json:"UDF02,omitempty"`
		UDF03 string  `json:"UDF03,omitempty"`
		UDF04 string  `json:"UDF04,omitempty"`
		UDF05 string  `json:"UDF05,omitempty"`
		UDF06 float64 `json:"UDF06,omitempty"`
		UDF07 float64 `json:"UDF07,omitempty"`
		UDF08 float64 `json:"UDF08,omitempty"`
		UDF09 float64 `json:"UDF09,omitempty"`
		UDF10 float64 `json:"UDF10,omitempty"`
	}{
		MF001: strings.TrimSpace(m.MF001),
		MF002: strings.TrimSpace(m.MF002),
		MF003: strings.TrimSpace(m.MF003),
		MF004: strings.TrimSpace(m.MF004),
		MF005: strings.TrimSpace(m.MF005),
		MF006: strings.TrimSpace(m.MF006),
		MF007: strings.TrimSpace(m.MF007),
		MF008: m.MF008,
		MF009: m.MF009,
		MF010: strings.TrimSpace(m.MF010),
		MF011: strings.TrimSpace(m.MF011),
		MF012: m.MF012,
		MF013: strings.TrimSpace(m.MF013),
		MF014: m.MF014,
		MF015: m.MF015,
		MF016: strings.TrimSpace(m.MF016),
		MF017: strings.TrimSpace(m.MF017),
		MF018: strings.TrimSpace(m.MF018),
		MF019: strings.TrimSpace(m.MF019),
		MF020: strings.TrimSpace(m.MF020),
		MF021: m.MF021,
		MF022: m.MF022,
		MF023: strings.TrimSpace(m.MF023),
		MF024: strings.TrimSpace(m.MF024),
		MF025: strings.TrimSpace(m.MF025),
		MF026: strings.TrimSpace(m.MF026),
		MF027: strings.TrimSpace(m.MF027),
		MF028: strings.TrimSpace(m.MF028),
		MF029: strings.TrimSpace(m.MF029),
		MF030: strings.TrimSpace(m.MF030),
		MF031: strings.TrimSpace(m.MF031),
		MF032: strings.TrimSpace(m.MF032),
		MF033: strings.TrimSpace(m.MF033),
		UDF01: strings.TrimSpace(m.UDF01),
		UDF02: strings.TrimSpace(m.UDF02),
		UDF03: strings.TrimSpace(m.UDF03),
		UDF04: strings.TrimSpace(m.UDF04),
		UDF05: strings.TrimSpace(m.UDF05),
		UDF06: m.UDF06,
		UDF07: m.UDF07,
		UDF08: m.UDF08,
		UDF09: m.UDF09,
		UDF10: m.UDF10,
	})
}

type SaleCopi04SearchReq struct {
	ME001 string `json:"ME001,omitempty"`
}
