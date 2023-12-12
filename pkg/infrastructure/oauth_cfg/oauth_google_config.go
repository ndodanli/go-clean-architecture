package oauthcfg

import (
	"github.com/ndodanli/go-clean-architecture/configs"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"strings"
)

var (
	GoogleOauth2Config *oauth2.Config
)

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func Init(googleOAuth2Cfg *configs.GoogleOauth2) {
	GoogleOauth2Config = &oauth2.Config{
		ClientID:     googleOAuth2Cfg.CLIENT_ID,
		ClientSecret: googleOAuth2Cfg.CLIENT_SECRET,
		//RedirectURL:  googleOAuth2Cfg.REDIRECT_URI,
		Scopes:   strings.Split(googleOAuth2Cfg.SCOPES, ","),
		Endpoint: google.Endpoint,
	}
}
