package models

// Settings
type UserSettings struct {
	Discoverable  bool   `json:"discoverable"`
	AllowMessages bool   `json:"allow_messages"`
	NotifyPush    bool   `json:"notify_push"`
	NotifyEmail   bool   `json:"notify_email"`
	UpdatedAt     string `json:"updated_at,omitempty"`
}

type UpdateSettingsRequest struct {
	Discoverable  *bool `json:"discoverable"`
	AllowMessages *bool `json:"allow_messages"`
	NotifyPush    *bool `json:"notify_push"`
	NotifyEmail   *bool `json:"notify_email"`
}
