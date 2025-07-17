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

	if m.config.DockerContainerName != "" {
		PrintColored("yellow", "Creating database dump via Docker...")
		// In Docker mode, we stream the output of pg_dump to a file on the host.
		args := []string{
			"exec", m.config.DockerContainerName,
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
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", m.config.HostPassword))

		dumpFile, err := os.Create(m.dumpFile)
		if err != nil {
			return fmt.Errorf("failed to create dump file on host: %v", err)
		}
		defer dumpFile.Close()
		cmd.Stdout = dumpFile
		cmd.Stderr = &stderr

	} else {
		PrintColored("yellow", "Creating database dump from HostDB...")
		// In host mode, pg_dump writes directly to the file.
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

// restoreDump restores a database dump, running psql on the host or in Docker.
func (m *Migrator) restoreDump() error {
	var cmd *exec.Cmd
	var stderr bytes.Buffer

	if m.config.DockerContainerName != "" {
		PrintColored("yellow", "Restoring dump to target via Docker...")
		// In Docker mode, we stream the dump file from the host to psql's stdin.
		args := []string{
			"exec", "-i", m.config.DockerContainerName,
			"psql",
			"-h", m.config.TargetHost,
			"-p", m.config.TargetPort,
			"-U", m.config.TargetUser,
			"-d", m.config.TargetDB,
			"--verbose",
		}
		cmd = exec.Command("docker", args...)
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", m.config.TargetPassword))

		// Redirect dump file content to stdin
		dumpFile, err := os.Open(m.dumpFile)
		if err != nil {
			return fmt.Errorf("failed to open dump file on host: %v", err)
		}
		defer dumpFile.Close()
		cmd.Stdin = dumpFile
		cmd.Stderr = &stderr

	} else {
		PrintColored("yellow", "Restoring dump to target PostgreSQL...")
		// In host mode, psql reads directly from the file.
		args := []string{
			"-h", m.config.TargetHost,
			"-p", m.config.TargetPort,
			"-U", m.config.TargetUser,
			"-d", m.config.TargetDB,
			"-f", m.dumpFile,
			"--verbose",
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
