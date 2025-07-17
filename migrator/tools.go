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
// or within the specified Docker container.
func (m *Migrator) checkPGTools() error {
	if m.config.DockerContainerName != "" {
		PrintColored("yellow", "Checking for Docker and PostgreSQL tools in container...")
		if _, err := exec.LookPath("docker"); err != nil {
			return fmt.Errorf("docker command not found. Please install Docker or run without Docker mode")
		}
		tools := []string{"pg_dump", "psql"}
		for _, tool := range tools {
			cmd := exec.Command("docker", "exec", m.config.DockerContainerName, "which", tool)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("%s not found in container %s. Please ensure PostgreSQL client tools are installed in the container", tool, m.config.DockerContainerName)
			}
		}
		PrintColored("green", "Docker and PostgreSQL tools found in container.")
		return nil
	}

	// Host mode: check for tools on the host machine
	PrintColored("yellow", "Checking PostgreSQL tools on host...")
	tools := []string{"pg_dump", "psql"}
	for _, tool := range tools {
		if _, err := exec.LookPath(tool); err != nil {
			return fmt.Errorf("%s not found. Please install PostgreSQL client tools", tool)
		}
	}
	PrintColored("green", "PostgreSQL tools found on host.")
	return nil
}
