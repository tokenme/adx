package common

type Stats struct {
	Pv     uint64  `json:"pv"`
	Uv     uint64  `json:"uv"`
	Clicks uint64  `json:"clicks"`
	Ctr    float64 `json:"ctr"`
}

type DateStats struct {
	Stats
	Date string `json:"date"`
}

type HourStats struct {
	Stats
	Hour uint `json:"hour"`
}

type CountryStats struct {
	Stats
	Name string `json:"name"`
	Id   uint   `json:"id"`
}

type BrowserTypeStats struct {
	Stats
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type OsStats struct {
	Stats
	Name string `json:"name"`
}
