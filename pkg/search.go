// Package api
// This file is part of Go Forensics (https://www.goforensics.io/)
// Copyright (C) 2022 Marten Mooij (https://www.mooijtech.com/)
package api

import (
	"encoding/json"
	core "github.com/mooijtech/goforensics-core/pkg"
	"net/http"
)

// handleSearch handles the search endpoint.
func (server *Server) handleSearch() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method == "POST" {
			_, project, err := server.AuthenticateRequest(request)

			if err != nil {
				Logger.Errorf("Failed to authenticate request: %s", err)
				http.Error(responseWriter, "Failed to authenticate request.", http.StatusBadRequest)
				return
			}

			var requestBody map[string]interface{}

			err = json.NewDecoder(request.Body).Decode(&requestBody)

			if err != nil {
				Logger.Errorf("Failed to decode request body: %s", err)
				http.Error(responseWriter, "Failed to decode request body.", http.StatusBadRequest)
				return
			}

			requestTreeNodeUUIDs := requestBody["treeNodeUUIDs"] // TODO - Rename this to folderUUIDs
			requestMessageUUID := requestBody["messageUUID"]
			requestQuery := requestBody["query"]

			Logger.Infof("Request tree node UUIDs: %s", requestTreeNodeUUIDs)

			if requestTreeNodeUUIDs != nil && len(requestTreeNodeUUIDs.([]interface{})) != 0 {
				// Get messages from the specified folders.
				var treeNodeUUIDs []string

				// Get all the children of this folder.
				for _, treeNodeUUID := range requestTreeNodeUUIDs.([]interface{}) {
					treeNodeUUIDs = append(treeNodeUUIDs, treeNodeUUID.(string))

					treeNodeChildrenUUIDs, err := core.WalkTreeNodeChildrenUUIDs(treeNodeUUID.(string), project)

					if err != nil {
						Logger.Errorf("Failed to get tree node children UUIDs: %s", err)
						http.Error(responseWriter, "Failed to get tree node children UUIDs.", http.StatusInternalServerError)
						return
					}

					treeNodeUUIDs = append(treeNodeUUIDs, treeNodeChildrenUUIDs...)
				}

				Logger.Infof("Tree node UUIDs (with children): %s", treeNodeUUIDs)

				messageRows, err := core.GetMessagesFromFolders(treeNodeUUIDs, project)

				if err != nil {
					Logger.Errorf("Failed to perform search: %s", err)
					http.Error(responseWriter, "Failed to perform search.", http.StatusInternalServerError)
					return
				}

				err = json.NewEncoder(responseWriter).Encode(&messageRows)

				if err != nil {
					Logger.Errorf("Failed to encode messages.")
					http.Error(responseWriter, "Failed to encode messages.", http.StatusInternalServerError)
					return
				}
			} else if requestMessageUUID != nil {
				var messageUUID string

				// Convert []interface{} to string
				for _, x := range requestMessageUUID.([]interface{}) {
					messageUUID = x.(string)
					break
				}

				message, err := core.GetMessageByUUID(messageUUID, project)

				if err != nil {
					Logger.Errorf("Failed to get message: %s", requestMessageUUID.(string))
					http.Error(responseWriter, "Failed to get message.", http.StatusInternalServerError)
					return
				}

				if err := json.NewEncoder(responseWriter).Encode(message); err != nil {
					Logger.Errorf("Failed to encode message: %s", err)
					http.Error(responseWriter, "Failed to encode message.", http.StatusInternalServerError)
					return
				}
			} else if requestQuery != nil {
				Logger.Infof("Request query: %s", requestQuery)

				messages, err := core.GetMessagesFromQuery(requestQuery.(string), project)

				if err != nil {
					Logger.Errorf("Failed to get messages from query: %s", err)
					http.Error(responseWriter, "Failed to get messages from query.", http.StatusInternalServerError)
					return
				}

				err = json.NewEncoder(responseWriter).Encode(messages)

				if err != nil {
					Logger.Errorf("Failed to write response: %s", err)
					http.Error(responseWriter, "Failed to write response.", http.StatusInternalServerError)
					return
				}
			}
		}
	}
}
