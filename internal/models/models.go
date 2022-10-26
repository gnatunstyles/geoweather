package models

type City struct {
	Name      string  `json:"name"`
	Country   string  `json:"country"`
	Lattitude float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

type PredictionParse struct {
	List []struct {
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
		Date string `json:"dt_txt"`
	} `json:"list"`
	City struct {
		Name string `json:"name"`
	} `json:"city"`
}

type Prediction struct {
	City string
	Temp float64
	Date string
}
