package area

type AreaRequest struct {
	NamaArea  string `json:"nama_area"`
	Kapasitas int    `json:"kapasitas"`
}

type AreaResponse struct {
	ID        uint   `json:"id"`
	NamaArea  string `json:"nama_area"`
	Kapasitas int    `json:"kapasitas"`
	Terisi    int    `json:"terisi"`
}