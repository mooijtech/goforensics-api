// Package api
// This file is part of Go Forensics (https://www.goforensics.io/)
// Copyright (C) 2022 Marten Mooij (https://www.mooijtech.com/)
package api

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/r3labs/sse/v2"
)

// Server represents our API server which communicates with users.
type Server struct {
	Router           *mux.Router
	Database         *sql.DB
	CookieStore      *sessions.CookieStore
	ServerSentEvents *sse.Server
}
