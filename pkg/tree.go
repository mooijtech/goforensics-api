// Package api
// This file is part of Go Forensics (https://www.goforensics.io/)
// Copyright (C) 2022 Marten Mooij (https://www.mooijtech.com/)
package api

import (
	"encoding/json"
	core "github.com/mooijtech/goforensics-core/pkg"
	"net/http"
)

// handleTree handles the tree endpoint.
func (server *Server) handleTree() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method == "GET" {
			_, project, err := server.AuthenticateRequest(request)

			if err != nil {
				Logger.Errorf("Failed to authenticate request: %s", err)
				http.Error(responseWriter, "Failed to authenticate request.", http.StatusUnauthorized)
				return
			}

			rootTreeNodes, err := core.GetRootTreeNodes(project.UUID, server.Database)

			if err != nil {
				Logger.Errorf("Failed to get root tree nodes by project UUID: %s", err)
				http.Error(responseWriter, "Failed to get root tree nodes by project UUID.", http.StatusInternalServerError)
				return
			}

			if len(rootTreeNodes) == 0 {
				Logger.Error("Failed to find root tree nodes by project UUID.")
				http.Error(responseWriter, "Failed to find root tree nodes by project UUID.", http.StatusInternalServerError)
				return
			}

			var treeNodeDTOs []core.TreeNodeDTO

			for i, treeNodeRoot := range rootTreeNodes {
				treeNodeDTOs = append(treeNodeDTOs, core.TreeNodeDTO{
					Value:    treeNodeRoot.FolderUUID,
					Label:    treeNodeRoot.Title,
					Children: []core.TreeNodeDTO{},
				})

				treeNodes, err := core.WalkTreeNodeChildren(treeNodeRoot.FolderUUID, project.UUID, server.Database)

				if err != nil {
					Logger.Errorf("Failed to walk tree node children: %s", err)
					http.Error(responseWriter, "Failed to walk tree node children.", http.StatusInternalServerError)
					return
				}

				for _, treeNode := range treeNodes {
					treeNodeDTOs[i].Children = append(treeNodeDTOs[i].Children, treeNode)
				}
			}

			if err := json.NewEncoder(responseWriter).Encode(&treeNodeDTOs); err != nil {
				Logger.Errorf("Failed to encode response: %s", err)
				http.Error(responseWriter, "Failed to encode response.", http.StatusInternalServerError)
				return
			}
		}
	}
}
