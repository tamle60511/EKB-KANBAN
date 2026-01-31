package dto

import (
	"encoding/json"
	"strings"
)

type SaleCOPMA struct {
	MA001  string `json:"MA001"`
	MA002  string `json:"MA002,omitempty"`
	MA004  string `json:"MA004,omitempty"`
	MA006  string `json:"MA006,omitempty"`
	MA015  string `json:"MA015,omitempty"`
	MA017  string `json:"MA017,omitempty"`
	MA076  string `json:"MA076,omitempty"`
	MA018  string `json:"MA018,omitempty"`
	MA019  string `json:"MA019,omitempty"`
	MA077  string `json:"MA077,omitempty"`
	MA002C string `json:"MA002C,omitempty"`
	MA015C string `json:"MA015C,omitempty"`
	MA017C string `json:"MA017C,omitempty"`
	MA018C string `json:"MA018C,omitempty"`
	MA019C string `json:"MA019C,omitempty"`
	MA077C string `json:"MA077C,omitempty"`
}

func (m SaleCOPMA) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		MA001  string `json:"MA001"`
		MA002  string `json:"MA002,omitempty"`
		MA004  string `json:"MA004,omitempty"`
		MA006  string `json:"MA006,omitempty"`
		MA015  string `json:"MA015,omitempty"`
		MA017  string `json:"MA017,omitempty"`
		MA076  string `json:"MA076,omitempty"`
		MA018  string `json:"MA018,omitempty"`
		MA019  string `json:"MA019,omitempty"`
		MA077  string `json:"MA077,omitempty"`
		MA002C string `json:"MA002C,omitempty"`
		MA015C string `json:"MA015C,omitempty"`
		MA017C string `json:"MA017C,omitempty"`
		MA018C string `json:"MA018C,omitempty"`
		MA019C string `json:"MA019C,omitempty"`
		MA077C string `json:"MA077C,omitempty"`
	}{
		MA001:  strings.TrimSpace(m.MA001),
		MA002:  strings.TrimSpace(m.MA002),
		MA004:  strings.TrimSpace(m.MA004),
		MA006:  strings.TrimSpace(m.MA006),
		MA015:  strings.TrimSpace(m.MA015),
		MA017:  strings.TrimSpace(m.MA017),
		MA076:  strings.TrimSpace(m.MA076),
		MA018:  strings.TrimSpace(m.MA018),
		MA019:  strings.TrimSpace(m.MA019),
		MA077:  strings.TrimSpace(m.MA077),
		MA002C: strings.TrimSpace(m.MA002C),
		MA015C: strings.TrimSpace(m.MA015C),
		MA017C: strings.TrimSpace(m.MA017C),
		MA018C: strings.TrimSpace(m.MA018C),
		MA019C: strings.TrimSpace(m.MA019C),
		MA077C: strings.TrimSpace(m.MA077C),
	})
}

type Types struct {
	MR002 string `json:"MR002"`
	MR003 string `json:"MR003"`
	MR004 string `json:"MR004"`
	MR005 string `json:"MR005"`
}

type SaleDepartment struct {
	ME001 string `json:"MC001"`
	ME002 string `json:"MC002"`
}

func (m SaleDepartment) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ME001 string `json:"MC001"`
		ME002 string `json:"MC002"`
	}{
		ME001: strings.TrimSpace(m.ME001),
		ME002: strings.TrimSpace(m.ME002),
	})
}

type SaleWorkshop struct {
	MB001 string `json:"MB001"`
	MB002 string `json:"MB002"`
}

func (m SaleWorkshop) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		MB001 string `json:"MB001"`
		MB002 string `json:"MB002"`
	}{
		MB001: strings.TrimSpace(m.MB001),
		MB002: strings.TrimSpace(m.MB002),
	})
}

type SaleItem struct {
	MB001 string `json:"MB001"`
	MB002 string `json:"MB002"`
	MB003 string `json:"MB003"`
	MB004 string `json:"MB004"`
	MB005 string `json:"MB005"`
	MB006 string `json:"MB006"`
	MB008 string `json:"MB008"`
	MB017 string `json:"MB017"`
	MC002 string `json:"MC002"`
}

func (m SaleItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		MB001 string `json:"MB001"`
		MB002 string `json:"MB002"`
		MB003 string `json:"MB003"`
		MB004 string `json:"MB004"`
		MB005 string `json:"MB005"`
		MB006 string `json:"MB006"`
		MB008 string `json:"MB008"`
		MB017 string `json:"MB017"`
		MC002 string `json:"MC002"`
	}{
		MB001: strings.TrimSpace(m.MB001),
		MB002: strings.TrimSpace(m.MB002),
		MB003: strings.TrimSpace(m.MB003),
		MB004: strings.TrimSpace(m.MB004),
		MB005: strings.TrimSpace(m.MB005),
		MB006: strings.TrimSpace(m.MB006),
		MB017: strings.TrimSpace(m.MB017),
		MC002: strings.TrimSpace(m.MC002),
	})
}

type SaleWarehouse struct {
	MC001 string `json:"MC001"`
	MC002 string `json:"MC002"`
	MC003 string `json:"MC003"`
}

func (m SaleWarehouse) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		MC001 string `json:"MC001"`
		MC002 string `json:"MC002"`
		MC003 string `json:"MC003"`
	}{
		MC001: strings.TrimSpace(m.MC001),
		MC002: strings.TrimSpace(m.MC002),
		MC003: strings.TrimSpace(m.MC003),
	})
}

type SaleMoney struct {
	MF001 string `json:"MF001"`
	MF002 string `json:"MF002"`
}

func (m SaleMoney) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		MF001 string `json:"MF001"`
		MF002 string `json:"MF002"`
	}{
		MF001: strings.TrimSpace(m.MF001),
		MF002: strings.TrimSpace(m.MF002),
	})
}
