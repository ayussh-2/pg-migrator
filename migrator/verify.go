package migrator

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func (m *Migrator) verifyMigration() error {
	PrintColored("yellow", "Verifying migration...")

	hostCount, err := m.getTableCount(
		m.config.HostHost, m.config.HostPort, m.config.HostDB,
		m.config.HostUser, m.config.HostPassword, m.config.HostSSLMode,
	)
	if err != nil {
		return fmt.Errorf("failed to get HostDB table count: %v", err)
	}

	targetCount, err := m.getTableCount(
		m.config.TargetHost, m.config.TargetPort, m.config.TargetDB,
		m.config.TargetUser, m.config.TargetPassword, m.config.TargetSSLMode,
	)
	if err != nil {
		return fmt.Errorf("failed to get target PostgreSQL table count: %v", err)
	}

	fmt.Printf("HostDB tables: %d\n", hostCount)
	fmt.Printf("Target PostgreSQL tables: %d\n", targetCount)

	if hostCount == targetCount {
		PrintColored("green", "Table count verification passed.")
	} else {
		PrintColored("yellow", "Warning: Table count mismatch. Please verify manually.")
	}

	return nil
}

func (m *Migrator) generatePgAdminInfo() error {
	PrintColored("yellow", "Generating pgAdmin connection information...")

	infoFile := filepath.Join(m.backupDir, "pgadmin_connection_info.txt")

	content := fmt.Sprintf(`pgAdmin Connection Information
==============================

Host: %s
Port: %s
Database: %s
Username: %s
Password: %s


Migration completed on: %s
Backup location: %s
`,
		m.config.TargetHost, m.config.TargetPort, m.config.TargetDB,
		m.config.TargetUser, m.config.TargetPassword,
		time.Now().Format("2006-01-02 15:04:05"),
		m.backupDir,
	)

	if err := os.WriteFile(infoFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write pgAdmin info file: %v", err)
	}

	PrintColored("green", fmt.Sprintf("pgAdmin connection info saved to: %s", infoFile))
	return nil
}
