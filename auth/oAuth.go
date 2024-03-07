package auth

import (
	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func configHolder() *oauth2.Config {
	clientID, clientSecret := oauthClientSecret()
	googleOauthConfig := &oauth2.Config{
		RedirectURL:  "http://localhost:8080/OauthCallback",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	return googleOauthConfig
}

func OAuthLoginURL() string {
	googleOauthConfig := configHolder()
	oauthStateString := oauthStateString()
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	return url
}

func CheckStateString(requestStateString string) bool {
	return requestStateString == oauthStateString()
}

func CodeExchange(code string) (*oauth2.Token, error) {
	googleOauthConfig := configHolder()
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	return token, err
}
