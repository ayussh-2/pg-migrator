package migrator

import (
	"fmt"
	"os/exec"
)

// PrintColored prints a message in the specified color.
func PrintColored(color, message string) {
	colors := map[string]string{
		"red":    "\033[0;31m",
		"green":  "\033[0;32m",
		"yellow": "\033[1;33m",
		"reset":  "\033[0m",
	}
	fmt.Printf("%s%s%s\n", colors[color], message, colors["reset"])
}

// checkPGTools checks if pg_dump and psql are available, either on the host
// or within the specified Docker containers.
func (m *Migrator) checkPGTools() error {
	// Check if Docker is available if any container is specified
	if m.config.HostDockerContainer != "" || m.config.TargetDockerContainer != "" {
		PrintColored("yellow", "Checking for Docker...")
		if _, err := exec.LookPath("docker"); err != nil {
			return fmt.Errorf("docker command not found. Please install Docker or run without Docker mode")
		}
	}

	// Check pg_dump availability
	if m.config.UseDockerForTools && m.config.HostDockerContainer != "" {
		PrintColored("yellow", "Checking pg_dump in Docker container...")
		cmd := exec.Command("docker", "exec", m.config.HostDockerContainer, "which", "pg_dump")
		if err := cmd.Run(); err != nil {
			PrintColored("yellow", "pg_dump not found in container, falling back to host...")
			// Fallback to host
			if _, err := exec.LookPath("pg_dump"); err != nil {
				return fmt.Errorf("pg_dump not found in container or on host. Please install PostgreSQL client tools")
			}
			m.config.UseDockerForTools = false
		} else {
			PrintColored("green", "pg_dump found in Docker container.")
		}
	} else {
		PrintColored("yellow", "Checking pg_dump on host machine...")
		if _, err := exec.LookPath("pg_dump"); err != nil {
			if m.config.TargetDockerContainer != "" {
				PrintColored("yellow", "pg_dump not found on host, trying Docker container...")
				cmd := exec.Command("docker", "exec", m.config.TargetDockerContainer, "which", "pg_dump")
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("pg_dump not found on host or in container. Please install PostgreSQL client tools")
				}
				m.config.UseDockerForTools = true
				m.config.HostDockerContainer = m.config.TargetDockerContainer
				PrintColored("green", "pg_dump found in Docker container, using Docker mode.")
			} else {
				return fmt.Errorf("pg_dump not found on host. Please install PostgreSQL client tools")
			}
		} else {
			PrintColored("green", "pg_dump found on host machine.")
		}
	}

	// Check psql availability for target
	if m.config.TargetDockerContainer != "" {
		PrintColored("yellow", "Checking psql in target Docker container...")
		cmd := exec.Command("docker", "exec", m.config.TargetDockerContainer, "which", "psql")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("psql not found in target container %s. Please ensure PostgreSQL client tools are installed", m.config.TargetDockerContainer)
		}
		PrintColored("green", "psql found in target container.")
	} else {
		PrintColored("yellow", "Checking psql on host machine...")
		if _, err := exec.LookPath("psql"); err != nil {
			return fmt.Errorf("psql not found on host. Please install PostgreSQL client tools")
		}
		PrintColored("green", "psql found on host machine.")
	}

	return nil
}
