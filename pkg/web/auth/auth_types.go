package auth

// Registration member request
type RegisterRequest struct {
	Otp   string `json:"otp" validate:"required"`
	Phone string `json:"phone" validate:"required"`
	Email string `json:"email" validate:"required"`
}

// Registration confirm member request
type RegisterConfirmRequest struct {
	Zonid           string `json:"zonid" validate:"required"`
	Name            string `json:"name" validate:"required"`
	Phone           string `json:"phone" validate:"required"`
	Email           string `json:"email" validate:"required"`
	Company         string `json:"company" validate:"required"`
	Brand           string `json:"brand" validate:"required"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
	DeviceId        string `json:"device_id"`
	IsH2H           bool   `json:"is_h2h"`
}

// Login member request
type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
