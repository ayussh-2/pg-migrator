package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	HostHost     string
	HostPort     string
	HostDB       string
	HostUser     string
	HostPassword string
	HostSSLMode  string
	HostDockerContainer string // Docker container for host DB operations

	TargetHost     string
	TargetPort     string
	TargetDB       string
	TargetUser     string
	TargetPassword string
	TargetSSLMode  string
	TargetDockerContainer string // Docker container for target DB operations
	
	// New field to control tool execution strategy
	UseDockerForTools bool // Whether to use Docker containers for pg_dump/psql
}

func LoadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Info: No .env file found, relying on environment variables.")
	}

	config := Config{
		// HostDB configuration
		HostHost:     os.Getenv("HOST_HOST"),
		HostPort:     os.Getenv("HOST_PORT"),
		HostDB:       os.Getenv("HOST_DB"),
		HostUser:     os.Getenv("HOST_USER"),
		HostPassword: os.Getenv("HOST_PASSWORD"),
		HostSSLMode:  os.Getenv("HOST_SSLMODE"),
		HostDockerContainer: os.Getenv("HOST_DOCKER_CONTAINER"),

		// Target PostgreSQL configuration
		TargetHost:     os.Getenv("TARGET_HOST"),
		TargetPort:     os.Getenv("TARGET_PORT"),
		TargetDB:       os.Getenv("TARGET_DB"),
		TargetUser:     os.Getenv("TARGET_USER"),
		TargetPassword: os.Getenv("TARGET_PASSWORD"),
		TargetSSLMode:  os.Getenv("TARGET_SSLMODE"),
		TargetDockerContainer: os.Getenv("TARGET_DOCKER_CONTAINER"),
		
		// Check if we should use Docker for tools
		UseDockerForTools: os.Getenv("USE_DOCKER_FOR_TOOLS") == "true",
	}

	// Auto-detect strategy: if target container is specified and no host container, use Docker for tools
	if config.TargetDockerContainer != "" && config.HostDockerContainer == "" {
		config.UseDockerForTools = true
		// Use target container for running pg_dump (with network access to external host)
		config.HostDockerContainer = config.TargetDockerContainer
	}

	if config.HostPort == "" {
		config.HostPort = "5432"
	}
	if config.TargetPort == "" {
		config.TargetPort = "5432"
	}
	if config.HostSSLMode == "" {
		config.HostSSLMode = "disable"
	}
	if config.TargetSSLMode == "" {
		config.TargetSSLMode = "disable"
	}

	missingVars := []string{}

	if config.HostHost == "" {
		missingVars = append(missingVars, "HOST_HOST")
	}
	if config.HostDB == "" {
		missingVars = append(missingVars, "HOST_DB")
	}
	if config.HostUser == "" {
		missingVars = append(missingVars, "HOST_USER")
	}
	if config.HostPassword == "" {
		missingVars = append(missingVars, "HOST_PASSWORD")
	}

	if config.TargetHost == "" {
		missingVars = append(missingVars, "TARGET_HOST")
	}
	if config.TargetDB == "" {
		missingVars = append(missingVars, "TARGET_DB")
	}
	if config.TargetUser == "" {
		missingVars = append(missingVars, "TARGET_USER")
	}
	if config.TargetPassword == "" {
		missingVars = append(missingVars, "TARGET_PASSWORD")
	}

	if len(missingVars) > 0 {
		return Config{}, fmt.Errorf("error: missing required environment variables: %s", strings.Join(missingVars, ", "))
	}

	return config, nil
}
