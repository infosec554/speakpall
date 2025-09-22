package models

import "time"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	ID           string     `json:"id"                      db:"id"`             // uuid
	Email        string     `json:"email"                   db:"email"`          // citext UNIQUE
	DisplayName  string     `json:"name"                    db:"display_name"`   // 1..80
	PasswordHash string     `json:"-"                       db:"password_hash"`  // jsonda chiqmaydi
	GoogleID     *string    `json:"google_id,omitempty"     db:"google_id"`      // NULLable UNIQUE
	AvatarURL    *string    `json:"avatar,omitempty"        db:"avatar_url"`     // NULLable
	Age          *int       `json:"age,omitempty"           db:"age"`            // NULL yoki 13..120
	Gender       *string    `json:"gender,omitempty"        db:"gender"`         // 'male' | 'female' | NULL
	CountryCode  *string    `json:"country_code,omitempty"  db:"country_code"`   // CHAR(2) yoki NULL
	TargetLang   *string    `json:"target_lang,omitempty"   db:"target_lang"`    // NULL
	Level        *string    `json:"level,omitempty"         db:"level"`          // NULL (A1..C2)
	Role         string     `json:"role"                    db:"role"`           // 'admin' | 'user' (DEFAULT 'user')
	CreatedAt    time.Time  `json:"created_at"              db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"              db:"updated_at"`
}

// Login uchun minimal ma'lumot
type LoginUser struct {
	ID           string `db:"id"`
	PasswordHash string `db:"password_hash"`
	Role         string `db:"role"`
}

// Signup
type SignupRequest struct {
	DisplayName string  `json:"name"          binding:"required,min=1,max=80"`
	Email       string  `json:"email"         binding:"required,email"`
	Password    string  `json:"password"      binding:"required,min=6"`
	TargetLang  *string `json:"target_lang,omitempty"`
	Level       *string `json:"level,omitempty"`
	CountryCode *string `json:"country_code,omitempty"`
}

type SignupResponse struct {
	ID string `json:"id"`
}

// Login
type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Role         string `json:"role"`
}

// Profil rasmini yangilash
type UpdateAvatarRequest struct {
	AvatarURL string `json:"avatar" binding:"required,url"`
}

// Parolni tiklash
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token          string `json:"token"            binding:"required"`
	NewPassword    string `json:"new_password"     binding:"required,min=6"`
	RepeatPassword string `json:"repeat_password"  binding:"required,eqfield=NewPassword"`
}

