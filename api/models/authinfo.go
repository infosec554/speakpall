package models



// AuthInfo — kontekstdagi autentifikatsiya ma’lumoti (masalan, middleware set qiladi).
type AuthInfo struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"` 
}

// LogoutRequest
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}



// Refresh tokenni yangilash
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Parolni almashtirish
type ChangePasswordRequest struct {
	OldPassword        string `json:"old_password"         binding:"required"`     
	NewPassword        string `json:"new_password"         binding:"required,min=6"` 

}
