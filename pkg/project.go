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
				http.Error(responseWriter, "Failed to authenticate user.", http.StatusBadRequest)
				return
			}

			projects, err := core.GetProjectsByUserUUID(user.Id, server.Database)

			if err != nil {
				Logger.Errorf("Failed to find projects from user UUID: %s", user.Id)
				http.Error(responseWriter, "Failed to find projects.", http.StatusNotFound)
				return
			}

			err = json.NewEncoder(responseWriter).Encode(&projects)

			if err != nil {
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

			var project core.Project

			err = json.NewDecoder(request.Body).Decode(&project)

			if err != nil {
				Logger.Errorf("Failed to decode request body: %s", err)
				http.Error(responseWriter, "Failed to decode request body.", http.StatusBadRequest)
				return
			}

			project.UUID = core.NewUUID()
			project.UserUUID = user.Id
			project.CreationDate = int(time.Now().Unix())

			Logger.Infof("Creating project (%s): %s...", project.UUID, project.Name)

			err = project.Save(server.Database)

			if err != nil {
				Logger.Errorf("Failed to save project: %s", err)
				http.Error(responseWriter, "Failed to save project.", http.StatusInternalServerError)
				return
			}

			directoryPaths := []string{
				core.GetProjectDirectory(project),
				core.GetProjectTempDirectory(project),
			}

			for _, directory := range directoryPaths {
				err = os.MkdirAll(directory, 0755)

				if err != nil {
					Logger.Errorf("Failed to create user project directories: %s", err)
					http.Error(responseWriter, "Failed to create user project directories.", http.StatusInternalServerError)
					return
				}
			}

			database, err := core.GetProjectDatabase(project)

			if err != nil {
				Logger.Error("Failed to get project database: %s", err)
				http.Error(responseWriter, "Failed to get project database.", http.StatusInternalServerError)
				return
			}

			err = core.CreateProjectDatabaseTables(database)

			if err != nil {
				Logger.Error("Failed to create project database tables: %s", err)
				http.Error(responseWriter, "Failed to create project database tables.", http.StatusInternalServerError)
				return
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

			err = json.NewEncoder(responseWriter).Encode(&project)

			if err != nil {
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
				http.Error(responseWriter, "Failed to authenticate user.", http.StatusBadRequest)
				return
			}

			var project core.Project

			if err := json.NewDecoder(request.Body).Decode(&project); err != nil {
				Logger.Errorf("Failed to decode request body: %s", err)
				http.Error(responseWriter, "Failed to decode request body.", http.StatusBadRequest)
				return
			}

			_, err = core.GetProjectByUUID(project.UUID, user.Id, server.Database)

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

			session.Values["projectUUID"] = project.UUID

			if err := session.Save(request, responseWriter); err != nil {
				Logger.Errorf("Failed to save session: %s", err)
				http.Error(responseWriter, "Failed to save session.", http.StatusInternalServerError)
				return
			}

			written, err := responseWriter.Write([]byte("{\"status\": \"OK\"}"))

			if err != nil {
				Logger.Errorf("Failed to write response (wrote %d bytes): %s", written, err)
				http.Error(responseWriter, "Failed to write response.", http.StatusInternalServerError)
				return
			}
		}
	}
}
