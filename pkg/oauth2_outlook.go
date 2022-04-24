// Package api
// This file is part of Go Forensics (https://www.goforensics.io/)
// Copyright (C) 2022 Marten Mooij (https://www.mooijtech.com/)
package api

import (
	"fmt"
	core "github.com/mooijtech/goforensics-core/pkg"
	"github.com/r3labs/sse/v2"
	"net/http"
)

// handleOutlookEmailsOAuth2 handles the Outlook OAuth2 redirect.
func (server *Server) handleOutlookEmailsOAuth2() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		http.Redirect(responseWriter, request, core.GetOutlookEmailsAuthURL(), http.StatusTemporaryRedirect)
	}
}

// handleOutlookEmailsOAuth2Callback handles the Outlook OAuth2 callback.
func (server *Server) handleOutlookEmailsOAuth2Callback() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		user, project, err := server.AuthenticateRequest(request)

		if err != nil {
			Logger.Errorf("Failed to authenticate request: %s", err)
			http.Error(responseWriter, "Failed to authenticate request.", http.StatusBadRequest)
			return
		}

		token, err := core.GetOutlookEmailsAccessToken(request)

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
			if err := core.ParseOutlookIMAPEmails(project, user.Identity.Traits.(map[string]interface{})["imapEmail"].(string), token, &progressPercentageChannel); err != nil {
				Logger.Errorf("Failed to parse Outlook IMAP emails: %s", err)
				return
			}
		}()

		http.Redirect(responseWriter, request, "http://localhost:3000/loading?projectId="+project.UUID, http.StatusTemporaryRedirect)
	}
}

// handleOutlookEmailsOAuth2 handles the Outlook OAuth2 redirect.
func (server *Server) handleOutlookUserProfileOAuth2() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		http.Redirect(responseWriter, request, core.GetOutlookUserProfileAuthURL(), http.StatusTemporaryRedirect)
	}
}

// handleOutlookEmailsOAuth2Callback handles the Outlook OAuth2 callback.
func (server *Server) handleOutlookUserProfileOAuth2Callback() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		token, err := core.GetOutlookUserProfileAccessToken(request)

		if err != nil {
			Logger.Errorf("Failed to get Outlook user profile OAuth2 token: %s", err)
			http.Error(responseWriter, "Failed to get Outlook user profile OAuth2 token.", http.StatusInternalServerError)
			return
		}

		userEmail, err := core.GetOutlookUserProfile(token)

		if err != nil {
			Logger.Errorf("Failed to get Outlook user profile: %s", err)
			http.Error(responseWriter, "Failed to get Outlook user profile.", http.StatusInternalServerError)
			return
		}

		// TODO - Update the user traits with the userEmail (for IMAP -> imapEmail trait).
		Logger.Infof("Got user email: %s", userEmail)

		http.Redirect(responseWriter, request, "http://localhost:1337/outlook/emails/auth", http.StatusTemporaryRedirect)
	}
}

func (server *Server) handleOutlookLoading() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		server.ServerSentEvents.ServeHTTP(responseWriter, request)
	}
}
