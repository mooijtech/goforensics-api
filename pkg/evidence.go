// Package api
// This file is part of Go Forensics (https://www.goforensics.io/)
// Copyright (C) 2022 Marten Mooij (https://www.mooijtech.com/)
package api

import (
	"encoding/json"
	core "github.com/mooijtech/goforensics-core/pkg"
	"net/http"
)

// handleEvidence handles the evidence endpoint.
func (server *Server) handleEvidence() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		_, project, err := server.AuthenticateRequest(request)

		if err != nil {
			Logger.Errorf("Failed to authenticate request: %s", err)
			http.Error(responseWriter, "Failed to authenticate request.", http.StatusUnauthorized)
			return
		}

		if request.Method == "POST" {
			var requestMap map[string]string

			if err := json.NewDecoder(request.Body).Decode(&requestMap); err != nil {
				Logger.Errorf("Failed to decode request body: %s", err)
				http.Error(responseWriter, "Failed to decode request body.", http.StatusBadRequest)
				return
			}

			evidenceFileName, ok := requestMap["fileName"]

			if !ok {
				Logger.Errorf("Failed to get evidence file name.")
				http.Error(responseWriter, "Failed to get evidence file name.", http.StatusBadRequest)
				return
			}

			evidenceFileHash, ok := requestMap["fileHash"]

			if !ok {
				Logger.Errorf("Failed to get evidence file hash.")
				http.Error(responseWriter, "Failed to get evidence file hash.", http.StatusBadRequest)
				return
			}

			var evidence core.Evidence

			evidence.UUID = core.NewUUID()
			evidence.FileName = evidenceFileName
			evidence.FileHash = evidenceFileHash
			evidence.IsParsed = false

			if err := evidence.Save(server.Database); err != nil {
				Logger.Errorf("Failed to save evidence to database: %s", err)
				http.Error(responseWriter, "Failed to save evidence to database.", http.StatusInternalServerError)
				return
			}

			if err := core.AddProjectEvidence(project.UUID, evidence.UUID, server.Database); err != nil {
				Logger.Errorf("Failed to add project evidence: %s", err)
				http.Error(responseWriter, "Failed to add project evidence.", http.StatusInternalServerError)
				return
			}

			Logger.Infof("Indexing evidence (%s): %s...", evidence.FileHash, evidence.FileName)

			if err := evidence.Parse(project, server.Database); err != nil {
				Logger.Errorf("Failed to parse evidence: %s", err)
				http.Error(responseWriter, "Failed to parse evidence.", http.StatusInternalServerError)
				return
			}

			if _, err := responseWriter.Write([]byte("\"status\": \"OK\"}")); err != nil {
				Logger.Errorf("Failed to write response: %s", err)
				http.Error(responseWriter, "Failed to write response.", http.StatusInternalServerError)
				return
			}
		}
	}
}
