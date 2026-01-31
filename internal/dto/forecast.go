package dto

type Forecast struct {
	TD01 string `json:"TD01"`
	MKH  string `json:"MKH"`
	KH01 string `json:"KH01"`
	TD02 string `json:"TD02"`
	TD03 string `json:"TD03"`
	TD04 string `json:"TD04"`
	TD05 string `json:"TD05"`
	TD06 string `json:"TD06"`
	TD07 string `json:"TD07"`
	TD08 string `json:"TD08"`
	TD09 string `json:"TD09"`
	KH02 string `json:"KH02"`
	TD10 string `json:"TD10"`
	TD11 string `json:"TD11"`
	TD12 string `json:"TD12"`
	TD13 string `json:"TD13"`
	TD14 string `json:"TD14"`
	TD15 string `json:"TD15"`
}

type CombinedForecast struct {
	Columns []Columns     `json:"columns"`
	Details []SubForecast `json:"details"`
}

type Columns struct {
	TD01 string `json:"TD01"`
	MKH  string `json:"MKH"`
	KH01 string `json:"KH01"`
	TD02 string `json:"TD02"`
	TD03 string `json:"TD03"`
	TD04 string `json:"TD04"`
	TD05 string `json:"TD05"`
	TD06 string `json:"TD06"`
	TD07 string `json:"TD07"`
	TD08 string `json:"TD08"`
}

type SubForecast struct {
	TD09 string `json:"TD09"`
	KH02 string `json:"KH02"`
	TD10 string `json:"TD10"`
	TD11 string `json:"TD11"`
	TD12 string `json:"TD12"`
	TD13 string `json:"TD13"`
	TD14 string `json:"TD14"`
	TD15 string `json:"TD15"`
}
