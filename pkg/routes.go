// Package api
// This file is part of Go Forensics (https://www.goforensics.io/)
// Copyright (C) 2022 Marten Mooij (https://www.mooijtech.com/)
package api

import (
	"github.com/rs/cors"
	"net/http"
)

// Start registers our routes and starts the server.
func (server *Server) Start() {
	server.Router.Handle("/projects", server.handleProjects())
	server.Router.Handle("/setProject", server.handleSetProject())
	server.Router.Handle("/evidence", server.handleEvidence())
	server.Router.Handle("/tree", server.handleTree())
	server.Router.Handle("/search", server.handleSearch())
	server.Router.Handle("/bookmarks", server.handleBookmarks())
	server.Router.Handle("/bookmark", server.handleBookmark())
	server.Router.Handle("/bookmark/{bookmarkUUID}", server.handleBookmark())
	server.Router.Handle("/tags", server.handleTags())
	server.Router.Handle("/tag", server.handleTag())
	server.Router.Handle("/report", server.handleReport())
	server.Router.Handle("/export", server.handleExport())
	server.Router.Handle("/file/{userUUID}/{projectUUID}/{fileName}", server.handleFile())
	server.Router.Handle("/network", server.handleNetwork())
	server.Router.HandleFunc("/outlook/loading", server.handleOutlookLoading())

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000", "https://www.goforensics.io"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}).Handler(server.Router)

	Logger.Infof("Starting the API on http://127.0.0.1:1337")
	Logger.Fatal(http.ListenAndServe(":1337", corsHandler))
}
