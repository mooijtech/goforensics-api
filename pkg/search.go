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

// Constants defining the search types.
const (
	SearchTypeTree    = "TREE"
	SearchTypeMessage = "MESSAGE"
	SearchTypeQuery   = "QUERY"
)

// handleSearch handles the search endpoint.
func (server *Server) handleSearch() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method == "POST" {
			_, project, err := server.AuthenticateRequest(request)

			if err != nil {
				Logger.Errorf("Failed to authenticate request: %s", err)
				http.Error(responseWriter, "Failed to authenticate request.", http.StatusUnauthorized)
				return
			}

			var requestBody map[string]interface{}

			err = json.NewDecoder(request.Body).Decode(&requestBody)

			if err != nil {
				Logger.Errorf("Failed to decode request body: %s", err)
				http.Error(responseWriter, "Failed to decode request body.", http.StatusBadRequest)
				return
			}

			searchType := mux.Vars(request)["searchType"]

			switch searchType {
			case SearchTypeTree:
				// Get messages from the specified folders (tree nodes).
				requestTreeNodeUUIDs, ok := requestBody["treeNodeUUIDs"].([]interface{})

				if !ok || len(requestTreeNodeUUIDs) == 0 {
					Logger.Errorf("Failed to get request treeNodeUUIDs.")
					http.Error(responseWriter, "Failed to get request treeNodeUUIDs.", http.StatusBadRequest)
					return
				}

				var treeNodeUUIDs []string

				// Create the list of tree node UUIDs and walk the tree node children.
				for _, requestTreeNodeUUID := range requestTreeNodeUUIDs {
					treeNodeUUID, ok := requestTreeNodeUUID.(string)

					if !ok {
						Logger.Errorf("Failed to get request treeNodeUUID.")
						http.Error(responseWriter, "Failed to get request treeNodeUUID.", http.StatusBadRequest)
						return
					}

					treeNodeUUIDs = append(treeNodeUUIDs, treeNodeUUID)

					treeNodeChildrenUUIDs, err := core.WalkTreeNodeChildrenUUIDs(treeNodeUUID, project.UUID, server.Database)

					if err != nil {
						Logger.Errorf("Failed to get tree node children UUIDs: %s", err)
						http.Error(responseWriter, "Failed to get tree node children UUIDs.", http.StatusInternalServerError)
						return
					}

					treeNodeUUIDs = append(treeNodeUUIDs, treeNodeChildrenUUIDs...)
				}

				messages, err := core.GetMessagesFromFolders(treeNodeUUIDs, project.UUID, server.Database)

				if err != nil {
					Logger.Errorf("Failed to perform search: %s", err)
					http.Error(responseWriter, "Failed to perform search.", http.StatusInternalServerError)
					return
				}

				if err := json.NewEncoder(responseWriter).Encode(&messages); err != nil {
					Logger.Errorf("Failed to encode messages.")
					http.Error(responseWriter, "Failed to encode messages.", http.StatusInternalServerError)
					return
				}
			case SearchTypeMessage:
				// Get a specific message.
				messageUUID, ok := requestBody["messageUUID"].(string)

				if !ok {
					Logger.Errorf("Failed to get request messageUUID.")
					http.Error(responseWriter, "Failed to get request messageUUID.", http.StatusBadRequest)
					return
				}

				message, err := core.GetMessageByUUID(messageUUID, project.UUID, server.Database)

				if err != nil {
					Logger.Errorf("Failed to get message: %s", err)
					http.Error(responseWriter, "Failed to get message.", http.StatusInternalServerError)
					return
				}

				if err := json.NewEncoder(responseWriter).Encode(message); err != nil {
					Logger.Errorf("Failed to encode message: %s", err)
					http.Error(responseWriter, "Failed to encode message.", http.StatusInternalServerError)
					return
				}
			case SearchTypeQuery:
				// Search query.
				query, ok := requestBody["query"].(string)

				if !ok {
					Logger.Errorf("Failed to get search query.")
					http.Error(responseWriter, "Failed to get search query.", http.StatusBadRequest)
					return
				}

				messages, err := core.GetMessagesFromQuery(query, project.UUID, server.Database)

				if err != nil {
					Logger.Errorf("Failed to get messages from query: %s", err)
					http.Error(responseWriter, "Failed to get messages from query.", http.StatusInternalServerError)
					return
				}

				if err := json.NewEncoder(responseWriter).Encode(messages); err != nil {
					Logger.Errorf("Failed to write response: %s", err)
					http.Error(responseWriter, "Failed to write response.", http.StatusInternalServerError)
					return
				}
			}
		}
	}
}
