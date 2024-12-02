package api

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	Token string // base64 encoded value
}
