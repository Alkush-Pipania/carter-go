package authentication

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"-"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
