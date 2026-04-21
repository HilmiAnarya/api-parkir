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

type LoginResponse struct {
	Token       string `json:"token"`
	IDUser      uint   `json:"id_user"`
	NamaLengkap string `json:"nama_lengkap"`
	Role        string `json:"role"`
}

// Format standar balasan API (JSON Response)
type UserResponse struct {
	ID          uint   `json:"id"`
	NamaLengkap string `json:"nama_lengkap"`
	Username    string `json:"username"`
	Role        string `json:"role"`
	StatusAktif bool   `json:"status_aktif"`
}

type UpdateUserRequest struct {
	NamaLengkap string `json:"nama_lengkap"`
	Username    string `json:"username"`
	Password    string `json:"password"` // Opsional, dikosongkan jika tidak ingin ganti
	Role        string `json:"role"`
}