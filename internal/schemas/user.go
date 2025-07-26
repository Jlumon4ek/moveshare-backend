package schemas

type SignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Message string `json:"message"`
}

type SignInRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type SignInResponse struct {
	UserID       int64  `json:"user_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}
