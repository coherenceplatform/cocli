package cocli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"golang.org/x/oauth2"
)

// For dev (main/review)
const (
	clientID      = "O5AkI9iHd4Okb3DCmu1P0em4YXFjAPr5"
	authDomain    = "dev-mkiob4vl.us.auth0.com"
	credsFilename = "~/.cocli/.authtoken"
	// coherenceDomain = "aa-external-cocli.control-plane-review.coherence.coherencesites.com"
	coherenceDomain = "126bdeab-68f9-4d29-a22d-51f193623390-web.coherencedev.com"
	// coherenceDomain = "main.control-plane-review.coherence.coherencesites.com"
)

// TODO: For prod
// const (
// 	clientID      = "O5AkI9iHd4Okb3DCmu1P0em4YXFjAPr5"
// 	authDomain    = "dev-mkiob4vl.us.auth0.com"
// 	credsFilename = "~/.cocli/.authtoken"
//  coherenceDomain = "app.withcoherence.com"
// )

var oauthConfig = &oauth2.Config{
	ClientID: clientID,
	Endpoint: oauth2.Endpoint{
		AuthURL:       fmt.Sprintf("https://%s/authorize", authDomain),
		TokenURL:      fmt.Sprintf("https://%s/oauth/token", authDomain),
		DeviceAuthURL: fmt.Sprintf("https://%s/oauth/device/code", authDomain),
	},
	Scopes: []string{"offline_access", "openid", "email", "profile"},
}

func GetCliVersion() string {
	return "0.0.1"
}

func GetCoherenceApiPrefix() string {
	if strings.ToLower(os.Getenv("COHERENCE_DEV")) == "true" {
		return "/api/public/cli/v1"
	}

	return "/api/cli/v1"
}

func GetCoherenceDomain() string {
	return coherenceDomain
}

func GetAuthDomain() string {
	return authDomain
}

func CoherenceApiRequest(method string, url string, body io.Reader) (*http.Response, error) {
	return AuthenticatedRequest(
		method,
		url,
		body,
		GetRefreshedToken().IDToken,
	)
}

func OauthProviderApiRequest(method string, url string, body io.Reader) (*http.Response, error) {
	return AuthenticatedRequest(
		method,
		url,
		body,
		GetRefreshedToken().AccessToken,
	)
}

func AuthenticatedRequest(method string, url string, body io.Reader, bearer_token string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(
		method,
		url,
		body,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearer_token))

	res, err := client.Do(req)

	if res != nil && res.StatusCode == 401 {
		fmt.Println("Clearing credentials... login again with `cocli auth login` or replace COCLI_REFRESH_TOKEN with a new refresh token.")
		if CredsFileExists() {
			ClearCredentialsFile()
		}
	}

	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println(res.StatusCode)
		panic(FormatJSONOutput(bodyBytes))
	}

	return res, err
}

func GetRefreshedToken() *TokenWithIdToken {
	if !(CredsFileExists() || IsRefreshTokenVarSet()) {
		NotifyAuthRequired()
		panic("Unauthorized")
	}

	token := GetTokenFromCredsFile()
	if token == nil {
		token = &oauth2.Token{
			RefreshToken: os.Getenv("COCLI_REFRESH_TOKEN"),
		}
	}
	ctx := context.Background()
	tokenSource := &TokenWithIdTokenSource{
		TokenSource: oauthConfig.TokenSource(ctx, token),
	}
	// refreshes token automatically, only if needed
	refreshedToken, err := tokenSource.Token()
	if err != nil {
		fmt.Println("Error refreshing access token. Refresh token may be expired.")
		fmt.Println("Clearing credentials... login again with `cocli auth login` or replace COCLI_REFRESH_TOKEN with a new refresh token.")
		if CredsFileExists() {
			ClearCredentialsFile()
		}
		panic(err)
	}
	WriteTokenFile(refreshedToken)

	return refreshedToken
}

func GetOauthConfig() *oauth2.Config {
	return oauthConfig
}

func WriteTokenFile(token *TokenWithIdToken) {
	data, err := json.Marshal(token)
	if err != nil {
		panic(err)
	}

	filename := FormatFilename(credsFilename)
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		panic(err)
	}
}

func CredsFileExists() bool {
	filename := FormatFilename(credsFilename)

	_, err := os.Stat(filename)
	if err != nil {
		// File DNE
		return false
	}

	return true
}

func GetTokenFromCredsFile() *oauth2.Token {
	if !CredsFileExists() {
		return nil
	}

	filename := FormatFilename(credsFilename)

	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var token oauth2.Token
	err = json.Unmarshal(data, &token)
	if err != nil {
		panic(err)
	}

	return &token
}

func GetIdTokenFromCredsFile() *TokenWithIdToken {
	if !CredsFileExists() {
		return nil
	}

	filename := FormatFilename(credsFilename)

	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var token TokenWithIdToken
	err = json.Unmarshal(data, &token)
	if err != nil {
		panic(err)
	}

	return &token
}

func ClearCredentialsFile() bool {
	filename := FormatFilename(credsFilename)

	err := os.Remove(filename)
	if err != nil {
		fmt.Printf("Error clearing credentials: %v\n", err)
		return false
	}

	return true
}

func FormatFilename(filename string) string {
	if filename[:2] == "~/" {
		usr, err := user.Current()
		if err != nil {
			panic(err)
		}
		filename = filepath.Join(usr.HomeDir, filename[2:])
	}

	return filename
}

func IsRefreshTokenVarSet() bool {
	if refreshToken := os.Getenv("COCLI_REFRESH_TOKEN"); refreshToken == "" {
		return false
	}

	return true
}

func NotifyAuthRequired() {
	fmt.Print("Authentication required. Please set COCLI_REFRESH_TOKEN, or login with `cocli auth login`\n\n")
}

func FormatJSONOutput(bodyBytes []byte) string {
	var formattedRespBody bytes.Buffer
	err := json.Indent(&formattedRespBody, bodyBytes, "", "    ")
	if err != nil {
		panic(err)
	}

	return string(formattedRespBody.Bytes())
}
