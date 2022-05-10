// Package api
// This file is part of Go Forensics (https://www.goforensics.io/)
// Copyright (C) 2022 Marten Mooij (https://www.mooijtech.com/)
package api

import (
	"encoding/json"
	core "github.com/mooijtech/goforensics-core/pkg"
	"net/http"
	"os"
	"time"
)

// handleProjects handles the projects API endpoint.
func (server *Server) handleProjects() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method == "GET" {
			// Returns a list of all projects from the user.
			user, err := server.AuthenticateUser(request)

			if err != nil {
				Logger.Errorf("Failed to authenticate user: %s", err)
				http.Error(responseWriter, "Failed to authenticate user.", http.StatusUnauthorized)
				return
			}

			projects, err := core.GetProjectsByUser(user.Id, server.Database)

			if err != nil {
				Logger.Errorf("Failed to find projects from user: %s", err)
				http.Error(responseWriter, "Failed to find projects.", http.StatusNotFound)
				return
			}

			if err := json.NewEncoder(responseWriter).Encode(&projects); err != nil {
				Logger.Errorf("Failed to encode response: %s", err)
				http.Error(responseWriter, "Failed to encode response.", http.StatusInternalServerError)
				return
			}
		} else if request.Method == "POST" {
			// Creates a new project
			user, err := server.AuthenticateUser(request)

			if err != nil {
				Logger.Errorf("Failed to authenticate user: %s", err)
				http.Error(responseWriter, "Failed to authenticate user.", http.StatusBadRequest)
				return
			}

			var requestMap map[string]string

			if err := json.NewDecoder(request.Body).Decode(&requestMap); err != nil {
				Logger.Errorf("Failed to decode request body: %s", err)
				http.Error(responseWriter, "Failed to decode request body.", http.StatusBadRequest)
				return
			}

			projectName, ok := requestMap["projectName"]

			if !ok {
				Logger.Errorf("Failed to get projectName.")
				http.Error(responseWriter, "Failed to get projectName.", http.StatusBadRequest)
				return
			}

			var project core.Project

			project.UUID = core.NewUUID()
			project.Name = projectName
			project.CreationDate = int(time.Now().Unix())

			Logger.Infof("Creating project (%s): %s...", project.UUID, project.Name)

			if err := project.Save(server.Database); err != nil {
				Logger.Errorf("Failed to save project: %s", err)
				http.Error(responseWriter, "Failed to save project.", http.StatusInternalServerError)
				return
			}

			if err := core.AddProjectUser(project.UUID, user.Id, server.Database); err != nil {
				Logger.Errorf("Failed to add project to user: %s", err)
				http.Error(responseWriter, "Failed to add project to user.", http.StatusInternalServerError)
				return
			}

			directoryPaths := []string{
				core.GetProjectDirectory(project.UUID),
				core.GetProjectTempDirectory(project.UUID),
			}

			for _, directory := range directoryPaths {
				err = os.MkdirAll(directory, 0755)

				if err != nil {
					Logger.Errorf("Failed to create project directories: %s", err)
					http.Error(responseWriter, "Failed to create project directories.", http.StatusInternalServerError)
					return
				}
			}

			session, err := server.CookieStore.Get(request, "session")

			if err != nil {
				Logger.Errorf("Failed to get session: %s", err)
				http.Error(responseWriter, "Failed to get session.", http.StatusInternalServerError)
				return
			}

			session.Values["projectUUID"] = project.UUID

			if err := session.Save(request, responseWriter); err != nil {
				Logger.Errorf("Failed to save session: %s", err)
				http.Error(responseWriter, "Failed to save session.", http.StatusInternalServerError)
				return
			}

			if err := json.NewEncoder(responseWriter).Encode(&project); err != nil {
				Logger.Errorf("Failed to encode response: %s", err)
				http.Error(responseWriter, "Failed to encode response.", http.StatusInternalServerError)
				return
			}
		}
	}
}

// handleSetProject handle the setProject endpoint.
func (server *Server) handleSetProject() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method == "POST" {
			// Set the current project.
			user, err := server.AuthenticateUser(request)

			if err != nil {
				Logger.Errorf("Failed to authenticate user: %s", err)
				http.Error(responseWriter, "Failed to authenticate user.", http.StatusUnauthorized)
				return
			}

			var requestMap map[string]string

			if err := json.NewDecoder(request.Body).Decode(&requestMap); err != nil {
				Logger.Errorf("Failed to decode request body: %s", err)
				http.Error(responseWriter, "Failed to decode request body.", http.StatusBadRequest)
				return
			}

			projectUUID, ok := requestMap["projectUUID"]

			if !ok {
				Logger.Errorf("Failed to get projectUUID.")
				http.Error(responseWriter, "Failed to get projectUUID.", http.StatusBadRequest)
				return
			}

			if !core.ProjectHasUser(projectUUID, user.Id, server.Database) {
				Logger.Errorf("User is not assigned to this project.")
				http.Error(responseWriter, "User is not assigned to this project.", http.StatusBadRequest)
				return
			}

			_, err = core.GetProjectByUUID(projectUUID, server.Database)

			if err != nil {
				Logger.Errorf("Failed to get project by UUID: %s", err)
				http.Error(responseWriter, "Failed to get project by UUID.", http.StatusBadRequest)
				return
			}

			session, err := server.CookieStore.Get(request, "session")

			if err != nil {
				Logger.Errorf("Failed to get session: %s", err)
				http.Error(responseWriter, "Failed to get session.", http.StatusInternalServerError)
				return
			}

			session.Values["projectUUID"] = projectUUID

			if err := session.Save(request, responseWriter); err != nil {
				Logger.Errorf("Failed to save session: %s", err)
				http.Error(responseWriter, "Failed to save session.", http.StatusInternalServerError)
				return
			}

			if _, err := responseWriter.Write([]byte("{\"status\": \"OK\"}")); err != nil {
				Logger.Errorf("Failed to write response: %s", err)
				http.Error(responseWriter, "Failed to write response.", http.StatusInternalServerError)
				return
			}
		}
	}
}
