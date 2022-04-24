// Package api
// This file is part of Go Forensics (https://www.goforensics.io/)
// Copyright (C) 2022 Marten Mooij (https://www.mooijtech.com/)
package api

import (
	"encoding/json"
	core "github.com/mooijtech/goforensics-core/pkg"
	"net/http"
	"strings"
)

// handleExport handles the export endpoint.
func (server *Server) handleExport() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method == "POST" {
			_, project, err := server.AuthenticateRequest(request)

			if err != nil {
				Logger.Errorf("Failed to authenticate request: %s", err)
				http.Error(responseWriter, "Failed to authenticate request.", http.StatusBadRequest)
				return
			}

			var requestMap map[string]interface{}

			if err := json.NewDecoder(request.Body).Decode(&requestMap); err != nil {
				Logger.Errorf("Failed to decode request: %s", err)
				http.Error(responseWriter, "Failed to decode request.", http.StatusBadRequest)
				return
			}

			extensions, ok := requestMap["extensions"].(string)

			if !ok {
				Logger.Errorf("Failed to get type from request map: %s", requestMap["type"])
				http.Error(responseWriter, "Failed to get type from request map.", http.StatusBadRequest)
				return
			}

			exportPath, err := core.ExportAttachments(strings.Split(extensions, "\n"), project)

			if err != nil {
				Logger.Errorf("Failed to export attachments: %s", err)
				http.Error(responseWriter, "Failed to export attachments.", http.StatusInternalServerError)
				return
			}

			written, err := responseWriter.Write([]byte(exportPath))

			if err != nil {
				Logger.Errorf("Failed to write response (wrote %d bytes): %s", written, err)
				http.Error(responseWriter, "Failed to write response.", http.StatusInternalServerError)
				return
			}
		}
	}
}
