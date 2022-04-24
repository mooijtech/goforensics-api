// Package api
// This file is part of Go Forensics (https://www.goforensics.io/)
// Copyright (C) 2022 Marten Mooij (https://www.mooijtech.com/)
package api

import (
	"encoding/json"
	core "github.com/mooijtech/goforensics-core/pkg"
	"net/http"
)

// handleTags handles the tags' endpoint.
func (server *Server) handleTags() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method == "POST" {
			// Add tag.
			_, project, err := server.AuthenticateRequest(request)

			if err != nil {
				Logger.Errorf("Failed to authenticate request: %s", err)
				http.Error(responseWriter, "Failed to authenticate request.", http.StatusBadRequest)
				return
			}

			var tag core.Tag

			if err := json.NewDecoder(request.Body).Decode(&tag); err != nil {
				Logger.Errorf("Failed to decode request body: %s", err)
				http.Error(responseWriter, "Failed to decode request body.", http.StatusBadRequest)
				return
			}

			if err := tag.Save(project); err != nil {
				Logger.Errorf("Failed to save tag: %s", err)
				http.Error(responseWriter, "Failed to save tag.", http.StatusInternalServerError)
				return
			}

			if written, err := responseWriter.Write([]byte("{\"status\": \"OK\"}")); err != nil {
				Logger.Errorf("Failed to write response (wrote %d bytes): %s", written, err)
				http.Error(responseWriter, "Failed to write response.", http.StatusInternalServerError)
				return
			}
		}
	}
}

// handleTag handles the tag endpoint.
func (server *Server) handleTag() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method == "POST" {
			// Get tag.
			_, project, err := server.AuthenticateRequest(request)

			if err != nil {
				Logger.Errorf("Failed to authenticate request: %s", err)
				http.Error(responseWriter, "Failed to authenticate request.", http.StatusBadRequest)
				return
			}

			var requestMap map[string]string

			if err := json.NewDecoder(request.Body).Decode(&requestMap); err != nil {
				Logger.Errorf("Failed to decode request body: %s", err)
				http.Error(responseWriter, "Failed to decode request body.", http.StatusBadRequest)
				return
			}

			messageUUID, ok := requestMap["message_uuid"]

			if !ok {
				Logger.Errorf("Failed to get message UUID from request.")
				http.Error(responseWriter, "Failed to get message UUID from request.", http.StatusBadRequest)
				return
			}

			tag, err := core.GetTag(messageUUID, project)

			if err != nil {
				Logger.Errorf("Failed to get tag: %s", err)
				http.Error(responseWriter, "Failed to get tag.", http.StatusInternalServerError)
				return
			}

			if err := json.NewEncoder(responseWriter).Encode(tag); err != nil {
				Logger.Errorf("Failed to encode response: %s", err)
				http.Error(responseWriter, "Failed to encode response.", http.StatusInternalServerError)
				return
			}
		}
	}
}
