// api/schwab/auth.go
package schwab

import (
	"context"
	"fmt"
	"log"
	"net/http"

	//"os"
	"time"

	"golang.org/x/oauth2"
)

var (
	schwabAuthURL  = "https://api.schwabapi.com/v1/oauth/authorize"
	schwabTokenURL = "https://api.schwabapi.com/v1/oauth/token"
)

type OAuthClient struct {
	Config *oauth2.Config
	Token  *oauth2.Token
}

func NewAuthClient(clientId, clientSecret, redirectURL string) *OAuthClient {
	return &OAuthClient{
		Config: &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"read_only"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  schwabAuthURL,
				TokenURL: schwabTokenURL,
			},
		},
	}
}

func (c *OAuthClient) Authenticate() error {
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s", schwabAuthURL, c.Config.ClientID, c.Config.RedirectURL)
	fmt.Println("Visit URL:")
	fmt.Println(url)

	// make channel
	codeCh := make(chan string)
	// make http server
	srv := &http.Server{Addr: ":8080"}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		session := r.URL.Query().Get("session")
		fmt.Fprintf(w, "Authorization code received! You can close this window.\n")
		fmt.Printf("Code: %s\nSession: %s\n", code, session)

		codeCh <- code
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for code
	code := <-codeCh
	fmt.Printf("Received code from Schwab API: %s\n", code)

	// Shutdown server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)

	return nil
}
