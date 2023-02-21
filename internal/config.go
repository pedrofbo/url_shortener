package internal

import (
	"errors"
	"log"
	"os"
)

var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)

type Config struct {
	EntriesTableName        string
	ParticlesTableName      string
	BaseEndpoint            string
	ApiHandler              string
	DefaultRedirectEndpoint string
}

func LoadConfig() (*Config, error) {
	entriesTableName, ok := os.LookupEnv("ENTRIES_TABLE_NAME")
	if !ok {
		return nil, errors.New("Environment variable ENTRIES_TABLE_NAME is not set")
	}
	particlesTableName, ok := os.LookupEnv("PARTICLES_TABLE_NAME")
	if !ok {
		return nil, errors.New("Environment variable PARTICLES_TABLE_NAME is not set")
	}
	baseEndpoint, ok := os.LookupEnv("BASE_ENDPOINT")
	if !ok {
		baseEndpoint = "http://localhost:8080/"
	}
	handler, ok := os.LookupEnv("API_HANDLER")
	if !ok {
		handler = "GIN"
	}
	defaultRedirectEndpoint, ok := os.LookupEnv("DEFAULT_REDIRECT_ENDPOINT")
	if !ok {
		defaultRedirectEndpoint = "https://fun.pyoh.dev/"
	}
	config := Config{
		EntriesTableName:        entriesTableName,
		ParticlesTableName:      particlesTableName,
		BaseEndpoint:            baseEndpoint,
		ApiHandler:              handler,
		DefaultRedirectEndpoint: defaultRedirectEndpoint,
	}

	return &config, nil
}
