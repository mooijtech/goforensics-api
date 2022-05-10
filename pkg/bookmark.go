// Package api
// This file is part of Go Forensics (https://www.goforensics.io/)
// Copyright (C) 2022 Marten Mooij (https://www.mooijtech.com/)
package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	core "github.com/mooijtech/goforensics-core/pkg"
	"net/http"
)

// handleBookmarks handles the "bookmarks" endpoint.
func (server *Server) handleBookmarks() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		_, project, err := server.AuthenticateRequest(request)

		if err != nil {
			Logger.Errorf("Failed to authenticate request: %s", err)
			http.Error(responseWriter, "Failed to authenticate request.", http.StatusUnauthorized)
			return
		}

		if request.Method == "POST" {
			// Add bookmarks.
			var requestMap map[string]interface{}

			if err := json.NewDecoder(request.Body).Decode(&requestMap); err != nil {
				Logger.Errorf("Failed to decode request body: %s", err)
				http.Error(responseWriter, "Failed to decode request body.", http.StatusBadRequest)
				return
			}

			requestBookmarkUUIDs, ok := requestMap["bookmarks"].([]interface{})

			if !ok || len(requestBookmarkUUIDs) == 0 {
				Logger.Errorf("Failed to get request bookmark UUIDs.")
				http.Error(responseWriter, "Failed to get request bookmark UUIDs.", http.StatusBadRequest)
				return
			}

			for _, requestBookmarkUUID := range requestBookmarkUUIDs {
				if err := core.AddBookmark(requestBookmarkUUID.(string), project.UUID, server.Database); err != nil {
					Logger.Errorf("Failed to add bookmark: %s", err)
					http.Error(responseWriter, "Failed to add bookmark.", http.StatusInternalServerError)
					return
				}
			}

			if written, err := responseWriter.Write([]byte("{\"status\": \"OK\"}")); err != nil {
				Logger.Errorf("Failed to write response (wrote %d bytes): %s", written, err)
				http.Error(responseWriter, "Failed to write response.", http.StatusInternalServerError)
				return
			}
		} else if request.Method == "GET" {
			// Get bookmarks by project.
			bookmarks, err := core.GetBookmarksByProject(project.UUID, server.Database)

			if err != nil {
				Logger.Errorf("Failed to get bookmarks by project: %s", err)
				http.Error(responseWriter, "Failed to get projects by UUID.", http.StatusInternalServerError)
				return
			}

			if err := json.NewEncoder(responseWriter).Encode(bookmarks); err != nil {
				Logger.Errorf("Failed to write response: %s", err)
				http.Error(responseWriter, "Failed to write response.", http.StatusInternalServerError)
				return
			}
		}
	}
}

// handleBookmark handles the bookmark endpoint.
func (server *Server) handleBookmark() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		_, project, err := server.AuthenticateRequest(request)

		if err != nil {
			Logger.Errorf("Failed to authenticate request: %s", err)
			http.Error(responseWriter, "Failed to authenticate request.", http.StatusBadRequest)
			return
		}

		if request.Method == "DELETE" {
			// Delete bookmark.
			messageUUID := mux.Vars(request)["uuid"]

			if err := core.RemoveBookmark(messageUUID, project.UUID, server.Database); err != nil {
				Logger.Errorf("Failed to delete bookmark: %s", err)
				http.Error(responseWriter, "Failed to delete bookmark.", http.StatusInternalServerError)
				return
			}

			if _, err := responseWriter.Write([]byte("{\"status\": \"OK\"}")); err != nil {
				Logger.Errorf("Failed to write response: %s", err)
				http.Error(responseWriter, "Failed to write response.", http.StatusInternalServerError)
			}
		}
	}
}
