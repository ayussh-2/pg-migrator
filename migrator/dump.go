package migrator

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func (m *Migrator) createDump() error {
	// Create backup directory first
	if err := os.MkdirAll(m.backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %v", err)
	}

	var cmd *exec.Cmd
	var stderr bytes.Buffer

	if m.config.UseDockerForTools && m.config.HostDockerContainer != "" {
		PrintColored("yellow", "Creating database dump via Docker container with network access...")
		// Run pg_dump inside Docker container with host network access
		args := []string{
			"exec",
			"-e", fmt.Sprintf("PGPASSWORD=%s", m.config.HostPassword), // Pass password as environment variable
			m.config.HostDockerContainer,
			"pg_dump",
			"-h", m.config.HostHost,
			"-p", m.config.HostPort,
			"-U", m.config.HostUser,
			"-d", m.config.HostDB,
			"--verbose",
			"--no-owner",
			"--no-privileges",
			"--clean",
			"--if-exists",
		}
		cmd = exec.Command("docker", args...)

		// Create dump file on host and redirect Docker output to it
		dumpFile, err := os.Create(m.dumpFile)
		if err != nil {
			return fmt.Errorf("failed to create dump file on host: %v", err)
		}
		defer dumpFile.Close()
		cmd.Stdout = dumpFile
		cmd.Stderr = &stderr

	} else {
		PrintColored("yellow", "Creating database dump from host machine...")
		// Run pg_dump directly on host (for external DBs like Neon)
		args := []string{
			"-h", m.config.HostHost,
			"-p", m.config.HostPort,
			"-U", m.config.HostUser,
			"-d", m.config.HostDB,
			"--verbose",
			"--no-owner",
			"--no-privileges",
			"--clean",
			"--if-exists",
			"-f", m.dumpFile,
		}
		cmd = exec.Command("pg_dump", args...)
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", m.config.HostPassword))
		cmd.Stderr = &stderr
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pg_dump failed: %v\nOutput: %s", err, stderr.String())
	}

	if info, err := os.Stat(m.dumpFile); err == nil {
		PrintColored("green", fmt.Sprintf("Database dump created successfully: %s", m.dumpFile))
		PrintColored("green", fmt.Sprintf("Dump file size: %.2f MB", float64(info.Size())/1024/1024))
	}

	return nil
}

// restoreDump restores a database dump to the target database
func (m *Migrator) restoreDump() error {
	var cmd *exec.Cmd
	var stderr bytes.Buffer

	if m.config.TargetDockerContainer != "" {
		PrintColored("yellow", "Restoring dump to target via Docker container...")
		// Run psql inside the target Docker container
		args := []string{
			"exec", "-i",
			"-e", fmt.Sprintf("PGPASSWORD=%s", m.config.TargetPassword), // Pass password as environment variable
			m.config.TargetDockerContainer,
			"psql",
			"-h", m.config.TargetHost,
			"-p", m.config.TargetPort,
			"-U", m.config.TargetUser,
			"-d", m.config.TargetDB,
			"-v", "ON_ERROR_STOP=1", // Stop on first error
		}
		cmd = exec.Command("docker", args...)

		// Redirect dump file content to stdin
		dumpFile, err := os.Open(m.dumpFile)
		if err != nil {
			return fmt.Errorf("failed to open dump file on host: %v", err)
		}
		defer dumpFile.Close()
		cmd.Stdin = dumpFile
		cmd.Stderr = &stderr

	} else {
		PrintColored("yellow", "Restoring dump to target database...")
		// Run psql directly on host
		args := []string{
			"-h", m.config.TargetHost,
			"-p", m.config.TargetPort,
			"-U", m.config.TargetUser,
			"-d", m.config.TargetDB,
			"-f", m.dumpFile,
			"-v", "ON_ERROR_STOP=1", // Stop on first error
		}
		cmd = exec.Command("psql", args...)
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", m.config.TargetPassword))
		cmd.Stderr = &stderr
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("psql restore failed: %v\nOutput: %s", err, stderr.String())
	}

	PrintColored("green", "Database restore completed successfully.")
	return nil
}
