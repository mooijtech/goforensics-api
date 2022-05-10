// Package api
// This file is part of Go Forensics (https://www.goforensics.io/)
// Copyright (C) 2022 Marten Mooij (https://www.mooijtech.com/)
package api

import (
	core "github.com/mooijtech/goforensics-core/pkg"
	"net/http"
)

// handleReport handle the report endpoint.
func (server *Server) handleReport() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method == "POST" {
			_, project, err := server.AuthenticateRequest(request)

			if err != nil {
				Logger.Errorf("Failed to authenticate request: %s", err)
				http.Error(responseWriter, "Failed to authenticate request.", http.StatusUnauthorized)
				return
			}

			bookmarks, err := core.GetBookmarksByProject(project.UUID, server.Database)

			if err != nil {
				Logger.Errorf("Failed to get bookmarks by project: %s", err)
				http.Error(responseWriter, "Failed to get bookmarks by project.", http.StatusInternalServerError)
				return
			}

			outputPath, err := core.CreateHTMLReport(bookmarks, project)

			if err != nil {
				Logger.Errorf("Failed to create HTML report: %s", err)
				http.Error(responseWriter, "Failed to create HTML report.", http.StatusInternalServerError)
				return
			}

			if _, err := responseWriter.Write([]byte(outputPath)); err != nil {
				Logger.Errorf("Failed to write response: %s", err)
				http.Error(responseWriter, "Failed to write response.", http.StatusInternalServerError)
				return
			}
		}
	}
}
