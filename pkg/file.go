// Package api
// This file is part of Go Forensics (https://www.goforensics.io/)
// Copyright (C) 2022 Marten Mooij (https://www.mooijtech.com/)
package api

import (
	"fmt"
	"github.com/gorilla/mux"
	core "github.com/mooijtech/goforensics-core/pkg"
	"net/http"
)

// handleFile handles the file endpoint.
func (server *Server) handleFile() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method == "GET" {
			userUUID := mux.Vars(request)["userUUID"]
			projectUUID := mux.Vars(request)["projectUUID"]
			fileName := mux.Vars(request)["fileName"]

			// TODO - Verify the user exists and owns the project.

			responseWriter.Header().Set("Content-Type", "application/octet-stream")

			err := core.WriteFileToWriter(fmt.Sprintf("%s/%s/%s", userUUID, projectUUID, fileName), responseWriter)

			if err != nil {
				Logger.Errorf("Failed to write file to writer: %s", err)
				http.Error(responseWriter, "Failed to write file to writer.", http.StatusInternalServerError)
				return
			}
		}
	}
}
