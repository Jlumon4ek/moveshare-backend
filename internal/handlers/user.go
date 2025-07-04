package handlers

import (
	"encoding/json"
	"net/http"

	"moveshare/internal/service"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// SignUpRequest represents the sign-up request payload
type SignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignInRequest represents the sign-in request payload
type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignUpResponse represents the sign-up response
type SignUpResponse struct {
	Message string `json:"message"`
}

// SignInResponse represents the sign-in response
type SignInResponse struct {
	UserID       int64  `json:"user_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// SignUp handles user registration
// @Summary Register a new user
// @Description Creates a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param body body SignUpRequest true "User registration data"
// @Success 201 {object} SignUpResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /sign-up [post]
func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request payload"}`, http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, `{"error": "username, email, and password are required"}`, http.StatusBadRequest)
		return
	}

	err := h.userService.SignUp(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		http.Error(w, `{"error": "failed to create user"}`, http.StatusInternalServerError)
		return
	}

	resp := SignUpResponse{Message: "User created successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// SignIn handles user login
// @Summary Login a user
// @Description Authenticates a user and returns user details with access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param body body SignInRequest true "User login data"
// @Success 200 {object} SignInResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /sign-in [post]
func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request payload"}`, http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, `{"error": "email and password are required"}`, http.StatusBadRequest)
		return
	}

	user, accessToken, refreshToken, err := h.userService.SignIn(r.Context(), req.Email, req.Password)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "invalid password" {
			http.Error(w, `{"error": "invalid credentials"}`, http.StatusUnauthorized)
			return
		}
		http.Error(w, `{"error": "failed to authenticate"}`, http.StatusInternalServerError)
		return
	}

	resp := SignInResponse{
		UserID:       user.ID,
		Username:     user.Username,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
