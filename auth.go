package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

var oauthConfig *oauth2.Config

func openInBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // linux, freebsd, openbsd, netbsd
		cmd = "xdg-open"
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func handleAuth(profile string, profileConfig ProfileConfig) error {
	oauthConfig = &oauth2.Config{
		ClientID: profileConfig.OAuthClientID,
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/SASLogon/oauth/authorize", profileConfig.SASEndpoint),
			TokenURL: fmt.Sprintf("%s/SASlogon/oauth/token", profileConfig.SASEndpoint),
		},
		RedirectURL: "http://localhost:8080/callback",
		Scopes:      []string{"openid"},
	}

	authURL := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Println("Opening browser for authentication...")
	// Open the browser in various ways here at authURL
	if err := openInBrowser(authURL); err != nil {
		fmt.Printf("Please open the following URL manually:\n%s\n", authURL)
	}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Authorization code not found", http.StatusBadRequest)
			return
		}

		fmt.Fprintln(w, "Authentication Successful! You may close this tab.")

		token, err := oauthConfig.Exchange(context.Background(), code)
		if err != nil {
			log.Println("Failed to exchange token:", err)
			return
		}

		savedToken := &CredentialItem{
			AccessToken:  token.AccessToken,
			Expiry:       token.Expiry,
			RefreshToken: token.RefreshToken,
		}

		if err := SaveToken(profile, savedToken); err != nil {
			log.Println("Failed to save token:", err)
			return
		}

		log.Println("Authentication Complete!")
	})

	log.Println("Waiting for Authentication...")
	return http.ListenAndServe(":8080", nil)

}
