// Package api
// This file is part of Go Forensics (https://www.goforensics.io/)
// Copyright (C) 2022 Marten Mooij (https://www.mooijtech.com/)
package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	core "github.com/mooijtech/goforensics-core/pkg"
	"github.com/r3labs/sse/v2"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"net/url"
)

// GoForensicsAPIURL defines the URL to the API.
var GoForensicsAPIURL string

// GoForensicsDashboardURL defines the URL to the dashboard.
var GoForensicsDashboardURL string

// Variables defining our Microsoft secrets.
var (
	MicrosoftClientID     string
	MicrosoftClientSecret string
)

// init initializes our configuration variables.
func init() {
	for _, configurationVariable := range []string{
		"go_forensics_api_url", "go_forensics_dashboard_url",
		"microsoft_client_id", "microsoft_client_secret",
	} {
		if !viper.IsSet(configurationVariable) {
			Logger.Fatalf("unset %s configuration variable", configurationVariable)
		}
	}

	GoForensicsAPIURL = viper.GetString("go_forensics_api_url")
	GoForensicsDashboardURL = viper.GetString("go_forensics_dashboard_url")
	MicrosoftClientID = viper.GetString("microsoft_client_id")
	MicrosoftClientSecret = viper.GetString("microsoft_client_secret")
}

var MicrosoftProfileOAuth2Config = &oauth2.Config{
	ClientID:     MicrosoftClientID,
	ClientSecret: MicrosoftClientSecret,
	RedirectURL:  fmt.Sprintf("%s/microsoft/profile/callback", GoForensicsAPIURL),
	Scopes: []string{
		"User.Read",
		"https://graph.microsoft.com/User.Read",
	},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
		TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
	},
}

var MicrosoftEmailsOAuth2Config = &oauth2.Config{
	ClientID:     MicrosoftClientID,
	ClientSecret: MicrosoftClientSecret,
	RedirectURL:  fmt.Sprintf("%s/microsoft/emails/callback", GoForensicsAPIURL),
	Scopes: []string{
		"offline_access",
		"https://outlook.office.com/User.Read",
		"https://outlook.office.com/IMAP.AccessAsUser.All",
	},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
		TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
	},
}

// handleMicrosoftProfileOAuth2 handles the Outlook OAuth2 redirect.
func (server *Server) handleMicrosoftProfileOAuth2() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		http.Redirect(responseWriter, request, MicrosoftProfileOAuth2Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline), http.StatusTemporaryRedirect)
	}
}

// handleOutlookEmailsOAuth2Callback handles the Outlook OAuth2 callback.
func (server *Server) handleMicrosoftProfileOAuth2Callback() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		queryParts, err := url.ParseQuery(request.URL.RawQuery)

		if err != nil {
			Logger.Errorf("Failed to parse query: %s", err)
			http.Error(responseWriter, "Failed to parse query.", http.StatusInternalServerError)
			return
		}

		code := queryParts["code"][0]

		token, err := MicrosoftProfileOAuth2Config.Exchange(context.Background(), code)

		if err != nil {
			Logger.Errorf("Failed to get Outlook profile OAuth2 token: %s", err)
			http.Error(responseWriter, "Failed to get Outlook user profile OAuth2 token.", http.StatusInternalServerError)
			return
		}

		userEmail, err := getMicrosoftProfile(token.AccessToken)

		if err != nil {
			Logger.Errorf("Failed to get Outlook profile: %s", err)
			http.Error(responseWriter, "Failed to get Outlook profile.", http.StatusInternalServerError)
			return
		}

		Logger.Infof("Got user email: %s", userEmail)

		http.Redirect(responseWriter, request, fmt.Sprintf("%s/outlook/emails/auth", core.GoForensicsAPIURL), http.StatusTemporaryRedirect)
	}
}

// getMicrosoftProfile returns the user email.
func getMicrosoftProfile(token string) (string, error) {
	request, err := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me", nil)

	request.Header.Add("Authorization", "Bearer "+token)

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return "", err
	}

	defer func() {
		err := response.Body.Close()

		if err != nil {
			Logger.Errorf("Failed to close response body: %s", err)
		}
	}()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	var responseMap map[string]interface{}

	if err := json.Unmarshal(body, &responseMap); err != nil {
		return "", err
	}

	userEmail, ok := responseMap["userPrincipalName"].(string)

	if !ok {
		return "", errors.New("failed to get userPrincipalName from response")
	}

	return userEmail, nil
}

// handleMicrosoftEmailsOAuth2 handles the Outlook OAuth2 redirect.
func (server *Server) handleMicrosoftEmailsOAuth2() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		http.Redirect(responseWriter, request, MicrosoftEmailsOAuth2Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline), http.StatusTemporaryRedirect)
	}
}

// handleMicrosoftEmailsOAuth2Callback handles the Outlook OAuth2 callback.
func (server *Server) handleMicrosoftEmailsOAuth2Callback() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		user, project, err := server.AuthenticateRequest(request)

		if err != nil {
			Logger.Errorf("Failed to authenticate request: %s", err)
			http.Error(responseWriter, "Failed to authenticate request.", http.StatusBadRequest)
			return
		}

		queryParts, err := url.ParseQuery(request.URL.RawQuery)

		if err != nil {
			Logger.Errorf("Failed to parse query: %s", err)
			http.Error(responseWriter, "Failed to parse query.", http.StatusInternalServerError)
			return
		}

		code := queryParts["code"][0]

		token, err := MicrosoftEmailsOAuth2Config.Exchange(context.Background(), code)

		if err != nil {
			Logger.Errorf("Failed to get Outlook OAuth2 token: %s", err)
			http.Error(responseWriter, "Failed to get Outlook OAuth2 token.", http.StatusInternalServerError)
			return
		}

		progressPercentageChannel := make(chan int)

		go func() {
			server.ServerSentEvents.CreateStream(project.UUID)

			for percentage := range progressPercentageChannel {
				server.ServerSentEvents.Publish(project.UUID, &sse.Event{
					Data: []byte(fmt.Sprintf("%d", percentage)),
				})
			}

			server.ServerSentEvents.Publish(project.UUID, &sse.Event{
				Data: []byte(fmt.Sprintf("%d", -1)),
			})
		}()

		go func() {
			if err := core.ParseOutlookIMAPEmails(project, user.Identity.Traits.(map[string]interface{})["imapEmail"].(string), token.AccessToken, &progressPercentageChannel); err != nil {
				Logger.Errorf("Failed to parse Outlook IMAP emails: %s", err)
				return
			}
		}()

		http.Redirect(responseWriter, request, fmt.Sprintf("%s/loading?projectUUID=%s", GoForensicsDashboardURL, project.UUID), http.StatusTemporaryRedirect)
	}
}
