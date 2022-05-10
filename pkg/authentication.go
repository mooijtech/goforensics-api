// Package api
// This file is part of Go Forensics (https://www.goforensics.io/)
// Copyright (C) 2022 Marten Mooij (https://www.mooijtech.com/)
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	core "github.com/mooijtech/goforensics-core/pkg"
	"github.com/spf13/viper"
	"net/http"
)

// OryKratosURL defines the URL where Ory Kratos is running.
var OryKratosURL string

func init() {
	if !viper.IsSet("ory_kratos_url") {
		Logger.Fatalf("unset ory_kratos_url configuration variable")
	}

	OryKratosURL = viper.GetString("ory_kratos_url")
}

// AuthenticateRequest authenticates the request (user and current project).
func (server *Server) AuthenticateRequest(request *http.Request) (core.User, core.Project, error) {
	user, err := server.AuthenticateUser(request)

	if err != nil {
		return core.User{}, core.Project{}, errors.New("failed to get user by UUID")
	}

	session, err := server.CookieStore.Get(request, "session")

	if err != nil {
		return core.User{}, core.Project{}, err
	}

	projectUUID, ok := session.Values["projectUUID"].(string)

	if !ok {
		return core.User{}, core.Project{}, errors.New("projectUUID is not a string")
	}

	if !core.ProjectHasUser(projectUUID, user.Id, server.Database) {
		return core.User{}, core.Project{}, errors.New("project is not assigned to this user")
	}

	project, err := core.GetProjectByUUID(projectUUID, server.Database)

	if err != nil {
		return core.User{}, core.Project{}, err
	}

	return user, project, nil
}

// AuthenticateUser authenticates the user via Ory Kratos.
func (server *Server) AuthenticateUser(request *http.Request) (core.User, error) {
	cookie, err := request.Cookie("ory_kratos_session")

	if err != nil {
		return core.User{}, err
	}

	request, err = http.NewRequest("GET", fmt.Sprintf("%s/sessions/whoami", OryKratosURL), nil)

	if err != nil {
		return core.User{}, err
	}

	request.AddCookie(cookie)

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return core.User{}, err
	}

	defer func() {
		err := response.Body.Close()

		if err != nil {
			Logger.Errorf("Failed to close response body: %s", err)
		}
	}()

	var user core.User

	if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
		return core.User{}, err
	}

	return user, nil
}
