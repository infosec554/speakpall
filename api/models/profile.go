package models

// GET /user/me javobi uchun (asosiy profil)
type Profile struct {
	ID           string  `json:"id"`
	Email        string  `json:"email"`
	DisplayName  string  `json:"name"`
	AvatarURL    *string `json:"avatar,omitempty"`
	Age          *int    `json:"age,omitempty"`
	Gender       *string `json:"gender,omitempty"` // "male"|"female"|nil
	CountryCode  *string `json:"country_code,omitempty"`
	NativeLang   *string `json:"native_lang,omitempty"`
	TargetLang   *string `json:"target_lang,omitempty"`
	Level        *int    `json:"level,omitempty"` // 1..6
	About        *string `json:"about,omitempty"`
	Timezone     *string `json:"timezone,omitempty"`
	CreatedAt    string  `json:"created_at"`
}

// PATCH /user/me (qisman yangilash)
type UpdateProfileRequest struct {
	DisplayName *string `json:"name"          binding:"omitempty,min=1,max=80"`
	AvatarURL   *string `json:"avatar"        binding:"omitempty,url"`
	Age         *int    `json:"age"           binding:"omitempty,min=13,max=120"`
	Gender      *string `json:"gender"        binding:"omitempty,oneof=male female"`
	CountryCode *string `json:"country_code"  binding:"omitempty,len=2"`
	NativeLang  *string `json:"native_lang"   binding:"omitempty"`
	TargetLang  *string `json:"target_lang"   binding:"omitempty"`
	Level       *int    `json:"level"         binding:"omitempty,min=1,max=6"`
	About       *string `json:"about"         binding:"omitempty,max=2000"`
	Timezone    *string `json:"timezone"      binding:"omitempty"`
}

// Interests
type UpdateInterestsRequest struct {
	InterestIDs []int `json:"interest_ids" binding:"required,min=1,dive,gt=0"`
}
type InterestsResponse struct {
	InterestIDs []int `json:"interest_ids"`
}

// Match preferences
type MatchPreferences struct {
	TargetLang     *string  `json:"target_lang,omitempty"`
	MinLevel       *int     `json:"min_level,omitempty"` // 1..6
	MaxLevel       *int     `json:"max_level,omitempty"`
	GenderFilter   *string  `json:"gender_filter,omitempty"` // male|female|any
	MinRating      *int     `json:"min_rating,omitempty"`    // 1..5
	CountriesAllow []string `json:"countries_allow,omitempty"`
	UpdatedAt      string   `json:"updated_at,omitempty"`
}
type UpdateMatchPrefsRequest struct {
	TargetLang     *string  `json:"target_lang"`
	MinLevel       *int     `json:"min_level"       binding:"omitempty,min=1,max=6"`
	MaxLevel       *int     `json:"max_level"       binding:"omitempty,min=1,max=6"`
	GenderFilter   *string  `json:"gender_filter"   binding:"omitempty,oneof=male female any"`
	MinRating      *int     `json:"min_rating"      binding:"omitempty,min=1,max=5"`
	CountriesAllow []string `json:"countries_allow"`
}
