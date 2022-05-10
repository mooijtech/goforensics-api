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
			user, project, err := server.AuthenticateRequest(request)

			if err != nil {
				Logger.Errorf("Failed to authenticate request: %s", err)
				http.Error(responseWriter, "Failed to authenticate request.", http.StatusUnauthorized)
				return
			}

			fileName := mux.Vars(request)["fileName"]

			responseWriter.Header().Set("Content-Type", "application/octet-stream")

			if err := core.WriteFileToWriter(fmt.Sprintf("%s/%s/%s", user.Id, project.UUID, fileName), responseWriter); err != nil {
				Logger.Errorf("Failed to write file to writer: %s", err)
				http.Error(responseWriter, "Failed to write file to writer.", http.StatusInternalServerError)
				return
			}
		}
	}
}
