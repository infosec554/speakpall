package service

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2v2 "google.golang.org/api/oauth2/v2"

	"speakpall/api/models"
)

type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type GoogleService interface {
	ExchangeCodeForUser(ctx context.Context, code string) (*models.GoogleUser, error)
}

type googleService struct {
	oauthConfig *oauth2.Config
}

func NewGoogleService(config GoogleOAuthConfig) GoogleService {
	return &googleService{
		oauthConfig: &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			RedirectURL:  config.RedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (g *googleService) ExchangeCodeForUser(ctx context.Context, code string) (*models.GoogleUser, error) {
	tok, err := g.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("oauth exchange failed: %w", err)
	}

	oas, err := oauth2v2.New(g.oauthConfig.Client(ctx, tok))
	if err != nil {
		return nil, fmt.Errorf("oauth client init failed: %w", err)
	}

	ui, err := oas.Userinfo.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("userinfo fetch failed: %w", err)
	}



	return &models.GoogleUser{
		Email:    ui.Email,
		Name:     ui.Name,
		GoogleID: ui.Id,
		Picture:  ui.Picture, // models.GoogleUser da omitempty bor
	}, nil
}
