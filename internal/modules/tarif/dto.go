package tarif

type TarifRequest struct {
	JenisKendaraan string  `json:"jenis_kendaraan"` // motor, mobil, lainnya
	TarifPerJam    float64 `json:"tarif_per_jam"`
}