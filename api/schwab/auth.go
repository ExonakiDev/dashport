// api/schwab/auth.go
package schwab

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

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

// {
//
//	 "expires_in": 1800, //Number of seconds access_token is valid for
//	 "token_type": "Bearer",
//	 "scope": "api",
//	 "refresh_token": "{REFRESH_TOKEN_HERE}", //Valid for 7 days
//	 "access_token": "{ACCESS_TOKEN_HERE}", //Valid for 30 minutes
//	 "id_token": "{JWT_HERE}"
//	}

type TokenResponse struct {
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"Scope"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	IdToken      string `json:"id_token"`
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

func (c *OAuthClient) Authenticate() (string, error) {
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s", schwabAuthURL, c.Config.ClientID, c.Config.RedirectURL)
	fmt.Println("Visit URL:")
	fmt.Println(url)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Paste Code Returned in Browser URL:")
	code, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	code = strings.TrimSpace(code)
	fmt.Printf("Received code from Schwab API: %s\n", code)

	return code, err
}

// curl -X POST https://api.schwabapi.com/v1/oauth/token \
// -H 'Authorization: Basic {BASE64_ENCODED_Client_ID:Client_Secret} \
// -H 'Content-Type: application/x-www-form-urlencoded' \
// -d 'grant_type=authorization_code&code={AUTHORIZATION_CODE_VALUE}&redirect_uri=https://example_url.com/callback_example'

func (c *OAuthClient) GetToken(code string) oauth2.Token {
	log.Print("Attempting to get token..")
	clientIDSecret := c.Config.ClientID + ":" + c.Config.ClientSecret
	credentials := base64.StdEncoding.EncodeToString([]byte(clientIDSecret))

	data := url.Values{}
	data.Set("grant_type", fmt.Sprintf("authorization_code&code=%s&redirect_uri=%s", code, c.Config.RedirectURL))
	req, err := http.NewRequest("POST", schwabTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", credentials))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	dumpedRequest, err := httputil.DumpRequest(req, true)
	if err != nil {
		panic(err)
	}

	fmt.Println("Dumped Request")
	fmt.Printf("Sending Request: %s", string(dumpedRequest))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to create Request!")
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	var token oauth2.Token
	json.Unmarshal(body, &token)

	return token
}
