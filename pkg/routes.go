// Package api
// This file is part of Go Forensics (https://www.goforensics.io/)
// Copyright (C) 2022 Marten Mooij (https://www.mooijtech.com/)
package api

import (
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"net/http"
)

// AllowedOrigins defines the allowed CORS origins.
var AllowedOrigins []string

// init initializes the AllowedOrigins.
func init() {
	if !viper.IsSet("allowed_origins") {
		Logger.Fatalf("unset allowed_origins environment variable")
	}

	AllowedOrigins = viper.GetStringSlice("allowed_origins")
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
	server.Router.HandleFunc("/outlook/loading", server.handleOutlookLoading())

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   AllowedOrigins,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}).Handler(server.Router)

	Logger.Infof("Starting the API on http://127.0.0.1:1337")
	Logger.Fatal(http.ListenAndServe(":1337", corsHandler))
}
