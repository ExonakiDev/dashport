// api/schwab/auth.go
package schwab

import {
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
}

var {
	schwabAuthURL = "https://api.schwabapi.com/v1/oauth/authorize"
	schwabTokenURL = "https://api.schwabapi.com/v1/oauth/token"
}

type OAuthClient struct {
	Config *oauth2.Config
	Token *oauth2.Token
}

func NewAuthClient(clientId, clientSecret, redirectURL string) *OAuthClient{
	return &OAuthClient{
		Config: &oauth2.Config{
			ClientID: clientId,
			ClientSecret: clientSecret,
			RedirectURL: redirectURL,
			Scopes: []string{"read_only"},
			Endpoint: oauth2.Endpoint{
				AuthURL: schwabAuthURL,
				TokenURL: schwabTokenURL,
			},
		},
	}
}

func (c *OAuthClient) Authenticate() error {
	// TODO: Implement Authenticate
}
