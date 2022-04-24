// Package main
// This file is part of Go Forensics (https://www.goforensics.io/)
// Copyright (C) 2022 Marten Mooij (https://www.mooijtech.com/)
package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	api "github.com/mooijtech/goforensics-api/pkg"
	core "github.com/mooijtech/goforensics-core/pkg"
	"github.com/r3labs/sse/v2"
)

func main() {
	serverDatabase, err := core.GetServerDatabase()

	if err != nil {
		api.Logger.Fatalf("Failed to get database: %s", err)
		return
	}

	err = core.CreateServerDatabaseTables(serverDatabase)

	if err != nil {
		api.Logger.Fatalf("Failed to create server database tables: %s", err)
		return
	}

	server := api.Server{
		Router:           mux.NewRouter(),
		Database:         serverDatabase,
		CookieStore:      sessions.NewCookieStore(securecookie.GenerateRandomKey(32)),
		ServerSentEvents: sse.New(),
	}

	server.Start()
}
