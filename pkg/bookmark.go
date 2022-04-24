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

// handleBookmarks handles the bookmarks' endpoint.
func (server *Server) handleBookmarks() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		_, project, err := server.AuthenticateRequest(request)

		if err != nil {
			Logger.Errorf("Failed to authenticate request: %s", err)
			http.Error(responseWriter, "Failed to authenticate request.", http.StatusBadRequest)
			return
		}

		if request.Method == "POST" {
			// Add bookmark.
			var requestMap map[string]interface{}

			if err := json.NewDecoder(request.Body).Decode(&requestMap); err != nil {
				Logger.Errorf("Failed to decode request body: %s", err)
				http.Error(responseWriter, "Failed to decode request body.", http.StatusBadRequest)
				return
			}

			requestBookmarkUUIDs, ok := requestMap["bookmarks"].([]interface{})

			if !ok {
				Logger.Errorf("Failed to get request bookmark UUIDs: %s", requestMap["bookmarks"])
				http.Error(responseWriter, "Failed to get request bookmark UUIDs.", http.StatusBadRequest)
				return
			}

			for _, requestBookmarkUUID := range requestBookmarkUUIDs {
				var bookmark core.Bookmark

				bookmark.MessageUUID = requestBookmarkUUID.(string)

				if err := bookmark.Save(project); err != nil {
					Logger.Errorf("Failed to save bookmark: %s", err)
					http.Error(responseWriter, "Failed to save bookmark.", http.StatusInternalServerError)
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
			bookmarks, err := core.GetBookmarksByProject(project)

			if err != nil {
				Logger.Errorf("Failed to get bookmarks by project: %s", err)
				http.Error(responseWriter, "Failed to get projects by UUID.", http.StatusInternalServerError)
				return
			}

			var messages []core.Message

			for _, bookmark := range bookmarks {
				message, err := core.GetMessageByUUID(bookmark.MessageUUID, project)

				if err != nil {
					Logger.Errorf("Failed to get message by UUID: %s", err)
					http.Error(responseWriter, "Failed to get message by UUID.", http.StatusInternalServerError)
					return
				}

				messages = append(messages, message)
			}

			if err := json.NewEncoder(responseWriter).Encode(messages); err != nil {
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

		if request.Method == "POST" {
			// Get bookmark.
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

			bookmark, err := core.GetBookmark(messageUUID, project)

			if err != nil {
				Logger.Errorf("Failed to get bookmark: %s", err)
				http.Error(responseWriter, "Failed to get bookmark.", http.StatusInternalServerError)
				return
			}

			if err := json.NewEncoder(responseWriter).Encode(bookmark); err != nil {
				Logger.Errorf("Failed to encode response: %s", err)
				http.Error(responseWriter, "Failed to encode response.", http.StatusInternalServerError)
				return
			}
		} else if request.Method == "DELETE" {
			// Delete bookmark.
			bookmarkUUID := mux.Vars(request)["bookmarkUUID"]

			err := core.DeleteBookmark(bookmarkUUID, project)

			if err != nil {
				Logger.Errorf("Failed to delete bookmark: %s", err)
				http.Error(responseWriter, "Failed to delete bookmark.", http.StatusInternalServerError)
				return
			}
		}
	}
}
