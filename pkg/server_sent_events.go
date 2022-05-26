// Package api
// This file is part of Go Forensics (https://www.goforensics.io/)
// Copyright (C) 2022 Marten Mooij (https://www.mooijtech.com/)
package api

import "net/http"

// handleLoading handles ServerSentEvents.
func (server *Server) handleLoading() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		server.ServerSentEvents.ServeHTTP(responseWriter, request)
	}
}
