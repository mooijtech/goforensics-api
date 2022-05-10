// Package api
// This file is part of Go Forensics (https://www.goforensics.io/)
// Copyright (C) 2022 Marten Mooij (https://www.mooijtech.com/)
package api

import (
	"encoding/json"
	core "github.com/mooijtech/goforensics-core/pkg"
	"net/http"
)

// handleNetwork handles the network endpoint.
func (server *Server) handleNetwork() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		_, project, err := server.AuthenticateRequest(request)

		if err != nil {
			Logger.Errorf("Failed to authenticate request: %s", err)
			http.Error(responseWriter, "Failed to authenticate request.", http.StatusUnauthorized)
			return
		}

		if request.Method == "GET" {
			network, err := core.GetNetwork(project.UUID, server.Database)

			if err != nil {
				Logger.Errorf("Failed to get network: %s", err)
				http.Error(responseWriter, "Failed to get network.", http.StatusInternalServerError)
				return
			}

			if err := json.NewEncoder(responseWriter).Encode(network); err != nil {
				Logger.Errorf("Failed to encode response: %s", err)
				http.Error(responseWriter, "Failed to encode response.", http.StatusInternalServerError)
				return
			}
		}
	}
}
