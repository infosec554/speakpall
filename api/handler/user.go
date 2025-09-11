package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"speakpall/api/models"
	"speakpall/pkg/jwt"
	"speakpall/pkg/password"
	"speakpall/pkg/security"
)

// SignUp godoc
// @Summary      Register a new user
// @Description  Register a new user (name, email, password)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user body models.SignupRequest true "Signup data"
// @Success      201 {object} models.SignupResponse
// @Failure      400 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /auth/signup [post]
func (h Handler) SignUp(c *gin.Context) {
	var req models.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Parolni hash qilish
	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		handleResponse(c, h.log, "failed to hash password", http.StatusInternalServerError, err.Error())
		return
	}
	req.Password = hashedPassword

	// Foydalanuvchini yaratish
	userID, err := h.services.User().Create(ctx, req)
	if err != nil {
		handleResponse(c, h.log, "failed to create user", http.StatusInternalServerError, err.Error())
		return
	}

	// UserID va xohlasangiz token qaytarishingiz mumkin
	handleResponse(c, h.log, "user created successfully", http.StatusCreated, models.SignupResponse{
		ID: userID,
	})
}

// Login godoc
// @Summary      User login
// @Description  User login with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login body models.LoginRequest true "Login credentials"
// @Success      200 {object} models.LoginResponse
// @Failure      400 {object} models.Response
// @Failure      401 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /auth/login [post]
// Login ...
func (h Handler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid login request", http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.services.User().GetForLoginByEmail(ctx, req.Email)
	if err != nil {
		handleResponse(c, h.log, "user not found", http.StatusUnauthorized, err.Error())
		return
	}
	if err := security.CompareHashAndPassword(user.Password, req.Password); err != nil {
		handleResponse(c, h.log, "invalid credentials", http.StatusUnauthorized, "email or password is incorrect")
		return
	}

	// ➜ faqat "role"
	at, err := jwt.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		handleResponse(c, h.log, "failed to generate access token", http.StatusInternalServerError, err.Error())
		return
	}
	rt, _, err := jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		handleResponse(c, h.log, "failed to generate refresh token", http.StatusInternalServerError, err.Error())
		return
	}

	resp := models.LoginResponse{
		ID:           user.ID,
		Role:         user.Role,
		AccessToken:  at,
		RefreshToken: rt,
	}
	handleResponse(c, h.log, "login successful", http.StatusOK, resp)
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Return new access & refresh token using a valid refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        refresh body models.RefreshTokenRequest true "Refresh token"
// @Success      200 {object} models.LoginResponse
// @Failure      400 {object} models.Response
// @Failure      401 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /auth/refresh-token [post]
func (h Handler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		handleResponse(c, h.log, "refresh_token is required", http.StatusBadRequest, nil)
		return
	}

	claims, err := jwt.ExtractClaims(req.RefreshToken)
	if err != nil {
		handleResponse(c, h.log, "invalid refresh token", http.StatusUnauthorized, err.Error())
		return
	}

	// typ tekshir (refresh bo‘lishi shart)
	if t, _ := claims["typ"].(string); t != "refresh" {
		handleResponse(c, h.log, "invalid token type", http.StatusUnauthorized, nil)
		return
	}

	userID := fmt.Sprint(claims["user_id"])
	role := fmt.Sprint(claims["role"]) // bo‘lmasa bo‘sh chiqmasligi uchun
	if userID == "" {
		handleResponse(c, h.log, "invalid claims in refresh token", http.StatusUnauthorized, nil)
		return
	}

	at, err := jwt.GenerateAccessToken(userID, role)
	if err != nil {
		handleResponse(c, h.log, "failed to generate access token", http.StatusInternalServerError, err.Error())
		return
	}
	rt, _, err := jwt.GenerateRefreshToken(userID)
	if err != nil {
		handleResponse(c, h.log, "failed to generate refresh token", http.StatusInternalServerError, err.Error())
		return
	}

	resp := models.LoginResponse{
		ID:           userID,
		Role:         role,
		AccessToken:  at,
		RefreshToken: rt,
	}
	handleResponse(c, h.log, "tokens refreshed", http.StatusOK, resp)
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change password (user must send old and new password)
// @Tags user
// @Accept json
// @Produce json
// @Param change_password body models.ChangePasswordRequest true "Change password"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /auth/change-password [post]
// @Security ApiKeyAuth
func (h Handler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		handleResponse(c, h.log, "unauthorized", http.StatusUnauthorized, nil)
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.services.User().ChangePassword(ctx, userID.(string), req.OldPassword, req.NewPassword)
	if err != nil {
		handleResponse(c, h.log, err.Error(), http.StatusBadRequest, nil)
		return
	}

	handleResponse(c, h.log, "password changed successfully", http.StatusOK, nil)
}

// @Summary      Google orqali login yoki registratsiya
// @Description  Google OAuth code orqali login yoki ro‘yxatdan o‘tish (JWT tokenlar qaytaradi)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data body models.GoogleAuthRequest true "Google authorization code"
// @Success      200 {object} models.LoginResponse
// @Failure      400 {object} models.Response
// @Failure      401 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /auth/google [post]
func (h Handler) GoogleAuth(c *gin.Context) {
	var req models.GoogleAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()

	googleUser, err := h.services.Google().ExchangeCodeForUser(ctx, req.Code)
	if err != nil {
		handleResponse(c, h.log, "Google login failed", http.StatusUnauthorized, err.Error())
		return
	}

	// Create or get user, returns userID (and ensures user exists)
	userID, err := h.services.User().GoogleAuth(ctx, googleUser.Email, googleUser.Name, googleUser.GoogleID)
	if err != nil {
		handleResponse(c, h.log, "failed to create/login user", http.StatusInternalServerError, err.Error())
		return
	}

	// Fetch role (so existing admins keep their role)
	u, err := h.services.User().GetByID(ctx, userID)
	if err != nil {
		handleResponse(c, h.log, "failed to load user", http.StatusInternalServerError, err.Error())
		return
	}
	role := u.Role
	if role == "" {
		role = "user"
	}

	// NEW helpers: access + refresh generated separately, with proper claims
	accessToken, err := jwt.GenerateAccessToken(userID, role)
	if err != nil {
		handleResponse(c, h.log, "failed to generate access token", http.StatusInternalServerError, err.Error())
		return
	}
	refreshToken, _, err := jwt.GenerateRefreshToken(userID)
	if err != nil {
		handleResponse(c, h.log, "failed to generate refresh token", http.StatusInternalServerError, err.Error())
		return
	}

	resp := models.LoginResponse{
		ID:           userID,
		Role:         role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	handleResponse(c, h.log, "login via google", http.StatusOK, resp)
}

func ExtractBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}

// Logout godoc
// @Summary      Logout (chiqish)
// @Description  JWT tokenlarni va sessionni bekor qiladi
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data body models.LogoutRequest false "Logout request (refresh_token optional)"
// @Success      200 {object} models.Response
// @Failure      401 {object} models.Response
// @Router       /auth/logout [post]
// @Security     ApiKeyAuth
func (h Handler) Logout(c *gin.Context) {
	accessToken := ExtractBearerToken(c)
	var req models.LogoutRequest
	_ = c.ShouldBindJSON(&req)

	// Contextni uzating!
	ctx := c.Request.Context()

	if accessToken != "" {
		_ = h.services.Redis().BlacklistToken(ctx, accessToken)
	}
	if req.RefreshToken != "" {
		_ = h.services.Redis().BlacklistToken(ctx, req.RefreshToken)
	}

	// Cookie ni tozalash (agar front uchun kerak bo‘lsa)
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	handleResponse(c, h.log, "Logged out successfully", http.StatusOK, nil)
}

// RequestPasswordReset godoc
// @Summary      Request password reset
// @Description  Send a password reset link to the user email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        email body models.PasswordResetRequest true "User email"
// @Success      200 {object} models.Response
// @Failure      400 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /auth/request-password-reset [post]
func (h Handler) RequestPasswordReset(c *gin.Context) {
	var req models.PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	// Foydalanuvchi uchun token yaratish va yuborish
	token, err := h.services.User().CreatePasswordResetToken(c.Request.Context(), req.Email)
	if err != nil {
		handleResponse(c, h.log, "failed to send reset link", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "password reset link sent", http.StatusOK, gin.H{"token": token})
}

// @Summary      Reset user password
// @Description  Reset the user password using the reset token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data body models.ResetPasswordRequest true "New password, repeat password, and token"
// @Success      200 {object} models.Response
// @Failure      400 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /auth/reset-password [post]
func (h *Handler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest

	// 1. Frontenddan kelgan yangi parol va takrorlashni tekshirish
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	// 2. Parollar mos kelmasligini tekshirish
	if req.NewPassword != req.RepeatPassword {
		handleResponse(c, h.log, "passwords do not match", http.StatusBadRequest, "Repeat password does not match the new password")
		return
	}

	// 3. Parol murakkabligini tekshirish
	err := password.ValidatePassword(req.NewPassword)
	if err != nil {
		handleResponse(c, h.log, "invalid password", http.StatusBadRequest, err.Error())
		return
	}

	// 3. Tokenni tasdiqlash va parolni yangilash
	userID, err := h.services.User().ValidatePasswordResetToken(c.Request.Context(), req.Token)
	if err != nil {
		handleResponse(c, h.log, "invalid or expired token", http.StatusUnauthorized, err.Error())
		return
	}

	// 4. Yangi parolni yangilash
	err = h.services.User().ResetPassword(c.Request.Context(), userID, req.NewPassword)
	if err != nil {
		handleResponse(c, h.log, "failed to reset password", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "password reset successfully", http.StatusOK, gin.H{"message": "Password has been successfully reset."})
}
