// Package api
// This file is part of Go Forensics (https://www.goforensics.io/)
// Copyright (C) 2022 Marten Mooij (https://www.mooijtech.com/)
package api

import "github.com/spf13/viper"

func init() {
	viper.SetConfigName("goforensics")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil {
		Logger.Fatalf("Failed to initialize configuration file: %s", err)
	}
}
