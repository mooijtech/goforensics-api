// Package api
// This file is part of Go Forensics (https://www.goforensics.io/)
// Copyright (C) 2022 Marten Mooij (https://www.mooijtech.com/)
package api

import (
	"encoding/json"
	core "github.com/mooijtech/goforensics-core/pkg"
	"net/http"
)

// handleEvidence handles the project API endpoint.
func (server *Server) handleEvidence() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		_, project, err := server.AuthenticateRequest(request)

		if err != nil {
			Logger.Errorf("Failed to authenticate request: %s", err)
			http.Error(responseWriter, "Failed to authenticate request.", http.StatusBadRequest)
			return
		}

		if request.Method == "POST" {
			var evidence core.Evidence

			err := json.NewDecoder(request.Body).Decode(&evidence)

			if err != nil {
				Logger.Errorf("Failed to decode request body: %s", err)
				http.Error(responseWriter, "Failed to decode request body.", http.StatusBadRequest)
				return
			}

			evidence.UUID = core.NewUUID()

			err = evidence.Save(project)

			if err != nil {
				Logger.Errorf("Failed to save evidence to database: %s", err)
				http.Error(responseWriter, "Failed to save evidence to database.", http.StatusInternalServerError)
				return
			}

			Logger.Infof("Indexing evidence (%s): %s...", evidence.FileName, evidence.FileName)

			err = evidence.Parse(project)

			if err != nil {
				Logger.Errorf("Failed to parse evidence: %s", err)
				http.Error(responseWriter, "Failed to parse evidence.", http.StatusInternalServerError)
				return
			}

			written, err := responseWriter.Write([]byte("\"status\": \"OK\"}"))

			if err != nil {
				Logger.Errorf("Failed to write response (%d written): %s", written, err)
				http.Error(responseWriter, "Failed to write response", http.StatusInternalServerError)
				return
			}
		}
	}
}
