// Package api
// This file is part of Go Forensics (https://www.goforensics.io/)
// Copyright (C) 2022 Marten Mooij (https://www.mooijtech.com/)
package api

import (
	"encoding/json"
	core "github.com/mooijtech/goforensics-core/pkg"
	"net/http"
)

// handleTags handles the "tags" endpoint.
func (server *Server) handleTags() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method == "POST" {
			// Add tag.
			_, project, err := server.AuthenticateRequest(request)

			if err != nil {
				Logger.Errorf("Failed to authenticate request: %s", err)
				http.Error(responseWriter, "Failed to authenticate request.", http.StatusUnauthorized)
				return
			}

			var requestMap map[string]interface{}

			if err := json.NewDecoder(request.Body).Decode(&requestMap); err != nil {
				Logger.Errorf("Failed to decode request body: %s", err)
				http.Error(responseWriter, "Failed to decode request body.", http.StatusBadRequest)
				return
			}

			tag, ok := requestMap["tag"].(string)

			if !ok {
				Logger.Errorf("Failed to get tag from request.")
				http.Error(responseWriter, "Failed to get tag from request.", http.StatusBadRequest)
				return
			}

			messageUUIDs, ok := requestMap["messageUUIDs"].([]interface{})

			if !ok {
				Logger.Errorf("Failed to get messageUUIDs from request.")
				http.Error(responseWriter, "Failed to get messageUUIDs from request.", http.StatusBadRequest)
				return
			}

			for _, messageUUID := range messageUUIDs {
				if err := core.AddTag(tag, messageUUID.(string), project.UUID, server.Database); err != nil {
					Logger.Errorf("Failed to add tag: %s", err)
					http.Error(responseWriter, "Failed to add tag.", http.StatusInternalServerError)
					return
				}
			}

			if _, err := responseWriter.Write([]byte("{\"status\": \"OK\"}")); err != nil {
				Logger.Errorf("Failed to write response: %s", err)
				http.Error(responseWriter, "Failed to write response.", http.StatusInternalServerError)
				return
			}
		}
	}
}
