package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ExonakiDev/dashport/api/schwab" // ← adjust this path to match your module name
	"golang.org/x/oauth2"
	"github.com/spf13/viper"
)

const (
	tokenFile    = "token.json"
)

func main() {
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Error reading in Config")
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Create OAuth client from your schwab package
	oauthClient := schwab.NewAuthClient(clientID, clientSecret, redirectURI)

	// Load existing token if available
	token, err := loadToken(tokenFile)
	if err != nil || !token.Valid() {
		fmt.Println("No valid token found — starting OAuth flow...")

		err = oauthClient.Authenticate()
		if err != nil {
			log.Fatalf("OAuth authentication failed: %v", err)
		}

		//saveToken(tokenFile, token)
	} else {
		fmt.Println("Loaded existing valid token.")
	}
}

func saveToken(path string, token *oauth2.Token) {
	file, err := os.Create(path)
	if err != nil {
		log.Printf("Failed to save token: %v", err)
		return
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(token); err != nil {
		log.Printf("Error encoding token: %v", err)
	}
	fmt.Printf("Token saved to %s (expires %s)\n", path, token.Expiry.Format(time.RFC822))
}

func loadToken(path string) (*oauth2.Token, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var token oauth2.Token
	if err := json.NewDecoder(file).Decode(&token); err != nil {
		return nil, err
	}
	return &token, nil
}

