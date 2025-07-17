package migrator

import (
	"fmt"
	"path/filepath"
	"pg-migrator/utils"
	"time"
)

type Migrator struct {
	config    utils.Config
	backupDir string
	dumpFile  string
}

func NewMigrator(config utils.Config) *Migrator {
	timestamp := time.Now().Format("20060102_150405")
	backupDir := fmt.Sprintf("./host_migration_%s", timestamp)
	dumpFile := filepath.Join(backupDir, "host_dump.sql")

	return &Migrator{
		config:    config,
		backupDir: backupDir,
		dumpFile:  dumpFile,
	}
}

func (m *Migrator) Migrate() error {
	PrintColored("green", "Starting HostDB to PostgreSQL Migration")
	fmt.Println("==================================================")

	steps := []struct {
		name string
		fn   func() error
	}{
		{"Check PostgreSQL tools", m.checkPGTools},
		{"Test database connections", m.testConnections},
		{"Create database dump", m.createDump},
		{"Restore database dump", m.restoreDump},
		{"Verify migration", m.verifyMigration},
		{"Generate pgAdmin info", m.generatePgAdminInfo},
	}

	for _, step := range steps {
		if err := step.fn(); err != nil {
			PrintColored("red", fmt.Sprintf("Error in %s: %v", step.name, err))
			return err
		}
	}

	fmt.Println()
	PrintColored("green", "Migration completed successfully!")
	fmt.Println("==================================================")
	fmt.Printf("Backup directory: %s\n", m.backupDir)
	fmt.Println("You can now connect to your PostgreSQL database using pgAdmin")
	fmt.Printf("Connection details are saved in: %s/pgadmin_connection_info.txt\n", m.backupDir)

	return nil
}
