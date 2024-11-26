package model

// Representasi model User di database.
type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Email        string `gorm:"not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"password_hash,omitempty"`
	Fullname     string `gorm:"not null" json:"fullname"` // Nama lengkap pengguna
	Address      string `json:"address,omitempty"`        // Alamat pengguna (opsional)
	Gender       string `json:"gender,omitempty"`         // Jenis kelamin pengguna (opsional)
	PhoneNumber  string `json:"phone_number,omitempty"`   // Nomor telepon pengguna (opsional)
}

// UserResponseDTO untuk data transfer object for ketika update profile.
type UserResponseDTO struct {
	ID          uint   `json:"id"`
	Email       string `json:"email"`
	Fullname    string `json:"fullname"`
	Address     string `json:"address,omitempty"`
	Gender      string `json:"gender,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

// UserRequestDTO untuk data transfer object for ketika update profile.
type UserRequestDTO struct {
	Email       string `json:"email" validate:"required,email"` // Email pengguna
	Password    string `json:"password" validate:"min=6"`       // Password pengguna
	Fullname    string `json:"fullname" validate:"required"`    // Nama lengkap pengguna
	Address     string `json:"address,omitempty"`               // Alamat pengguna (opsional)
	Gender      string `json:"gender,omitempty"`                // Jenis kelamin pengguna (opsional)
	PhoneNumber string `json:"phone_number,omitempty"`          // Nomor telepon pengguna (opsional)
}

// authenticationRequest mendefinisikan struktur permintaan untuk pendaftaran dan login.
type AuthenticationRequest struct {
	Email    string `json:"email" validate:"required,email"`    // Email harus berupa format email yang valid
	Password string `json:"password" validate:"required,min=6"` // Password harus memiliki panjang minimal 6 karakter
}
