package models

// GoogleAuthRequest — Google OAuth-dan olingan "authorization code"ni qabul qilish uchun
type GoogleAuthRequest struct {
	Code string `json:"code" binding:"required" example:"4/0AX4XfW..."` // OAuth authorization code
}

// GoogleUser — Google'dan keladigan foydalanuvchi ma'lumotlari
type GoogleUser struct {
	Email    string `json:"email"     example:"user@example.com"`
	Name     string `json:"name"      example:"John Doe"`
	GoogleID string `json:"google_id" example:"123456789012345678901"`
	Picture  string `json:"picture,omitempty" example:"https://lh3.googleusercontent.com/a-/AOh14Gg..."` // ixtiyoriy profil rasmi
}
