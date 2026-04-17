package user

// Request untuk menambah user baru atau register
type CreateUserRequest struct {
	NamaLengkap string `json:"nama_lengkap"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Role        string `json:"role"` // admin, petugas, owner
}

// Request untuk Login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Format standar balasan API (JSON Response)
type UserResponse struct {
	ID          uint   `json:"id"`
	NamaLengkap string `json:"nama_lengkap"`
	Username    string `json:"username"`
	Role        string `json:"role"`
	StatusAktif bool   `json:"status_aktif"`
}