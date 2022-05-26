// Package api
// This file is part of Go Forensics (https://www.goforensics.io/)
// Copyright (C) 2022 Marten Mooij (https://www.mooijtech.com/)
package api

import (
	"fmt"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"net/http"
)

// AllowedOrigins defines the allowed CORS origins.
var AllowedOrigins []string

// GoForensicsAPIPort defines the port the API will run on.
var GoForensicsAPIPort int

// init initializes the AllowedOrigins and GoForensicsAPIPort.
func init() {
	for _, configurationVariable := range []string{"allowed_origins", "go_forensics_api_port"} {
		if !viper.IsSet(configurationVariable) {
			Logger.Fatalf("unset %s configuration variable", configurationVariable)
		}
	}

	AllowedOrigins = viper.GetStringSlice("allowed_origins")
	GoForensicsAPIPort = viper.GetInt("go_forensics_api_port")
}

// Start registers our routes and starts the server.
func (server *Server) Start() {
	server.Router.Handle("/projects", server.handleProjects())
	server.Router.Handle("/setProject", server.handleSetProject())
	server.Router.Handle("/evidence", server.handleEvidence())
	server.Router.Handle("/tree", server.handleTree())
	server.Router.Handle("/search/{searchType}", server.handleSearch())
	server.Router.Handle("/bookmarks", server.handleBookmarks())
	server.Router.Handle("/bookmark/{uuid}", server.handleBookmark())
	server.Router.Handle("/tags", server.handleTags())
	server.Router.Handle("/report", server.handleReport())
	server.Router.Handle("/export", server.handleExport())
	server.Router.Handle("/file/{fileName}", server.handleFile())
	server.Router.Handle("/network", server.handleNetwork())
	server.Router.HandleFunc("/loading", server.handleLoading())

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   AllowedOrigins,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}).Handler(server.Router)

	Logger.Infof("Starting the API on http://0.0.0.0:%d", GoForensicsAPIPort)
	Logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", GoForensicsAPIPort), corsHandler))
}
