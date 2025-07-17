package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// database configuration
type Config struct {
	HostHost     string
	HostPort     string
	HostDB       string
	HostUser     string
	HostPassword string
	
	TargetHost     string
	TargetPort     string
	TargetDB       string
	TargetUser     string
	TargetPassword string
}

type Migrator struct {
	config    Config
	backupDir string
	dumpFile  string
}

func NewMigrator(config Config) *Migrator {
	timestamp := time.Now().Format("20060102_150405")
	backupDir := fmt.Sprintf("./host_migration_%s", timestamp)
	dumpFile := filepath.Join(backupDir, "host_dump.sql")
	
	return &Migrator{
		config:    config,
		backupDir: backupDir,
		dumpFile:  dumpFile,
	}
}

func printColored(color, message string) {
	colors := map[string]string{
		"red":    "\033[0;31m",
		"green":  "\033[0;32m",
		"yellow": "\033[1;33m",
		"reset":  "\033[0m",
	}
	fmt.Printf("%s%s%s\n", colors[color], message, colors["reset"])
}

func (m *Migrator) checkPGTools() error {
	printColored("yellow", "Checking PostgreSQL tools...")
	
	tools := []string{"pg_dump", "psql"}
	for _, tool := range tools {
		if _, err := exec.LookPath(tool); err != nil {
			return fmt.Errorf("%s not found. Please install PostgreSQL client tools", tool)
		}
	}
	
	printColored("green", "PostgreSQL tools found.")
	return nil
}

func (m *Migrator) testConnection(host, port, dbname, user, password string, name string) error {
	printColored("yellow", fmt.Sprintf("Testing %s connection...", name))
	
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %v", name, err)
	}
	defer db.Close()
	
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping %s: %v", name, err)
	}
	
	printColored("green", fmt.Sprintf("%s connection successful.", name))
	return nil
}

func (m *Migrator) testConnections() error {
	if err := m.testConnection(
		m.config.HostHost, m.config.HostPort, m.config.HostDB,
		m.config.HostUser, m.config.HostPassword, "HostDB",
	); err != nil {
		return err
	}
	
	if err := m.testConnection(
		m.config.TargetHost, m.config.TargetPort, m.config.TargetDB,
		m.config.TargetUser, m.config.TargetPassword, "Target PostgreSQL",
	); err != nil {
		return err
	}
	
	return nil
}

func (m *Migrator) createDump() error {
	printColored("yellow", "Creating database dump from HostDB...")
	
	// Create backup directory
	if err := os.MkdirAll(m.backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %v", err)
	}
	
	// Prepare pg_dump command
	cmd := exec.Command("pg_dump",
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
	)
	
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", m.config.HostPassword))
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pg_dump failed: %v\nOutput: %s", err, output)
	}
	
	if info, err := os.Stat(m.dumpFile); err == nil {
		printColored("green", fmt.Sprintf("Database dump created successfully: %s", m.dumpFile))
		printColored("green", fmt.Sprintf("Dump file size: %.2f MB", float64(info.Size())/1024/1024))
	}
	
	return nil
}

func (m *Migrator) restoreDump() error {
	printColored("yellow", "Restoring dump to target PostgreSQL...")
	
	cmd := exec.Command("psql",
		"-h", m.config.TargetHost,
		"-p", m.config.TargetPort,
		"-U", m.config.TargetUser,
		"-d", m.config.TargetDB,
		"-f", m.dumpFile,
		"--verbose",
	)
	
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", m.config.TargetPassword))
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("psql restore failed: %v\nOutput: %s", err, output)
	}
	
	printColored("green", "Database restore completed successfully.")
	return nil
}

func (m *Migrator) getTableCount(host, port, dbname, user, password string) (int, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return 0, err
	}
	defer db.Close()
	
	var count int
	query := "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'"
	err = db.QueryRow(query).Scan(&count)
	return count, err
}

func (m *Migrator) verifyMigration() error {
	printColored("yellow", "Verifying migration...")
	
	hostCount, err := m.getTableCount(
		m.config.HostHost, m.config.HostPort, m.config.HostDB,
		m.config.HostUser, m.config.HostPassword,
	)
	if err != nil {
		return fmt.Errorf("failed to get HostDB table count: %v", err)
	}
	
	// Get table count from target PostgreSQL
	targetCount, err := m.getTableCount(
		m.config.TargetHost, m.config.TargetPort, m.config.TargetDB,
		m.config.TargetUser, m.config.TargetPassword,
	)
	if err != nil {
		return fmt.Errorf("failed to get target PostgreSQL table count: %v", err)
	}
	
	fmt.Printf("HostDB tables: %d\n", hostCount)
	fmt.Printf("Target PostgreSQL tables: %d\n", targetCount)
	
	if hostCount == targetCount {
		printColored("green", "Table count verification passed.")
	} else {
		printColored("yellow", "Warning: Table count mismatch. Please verify manually.")
	}
	
	return nil
}

func (m *Migrator) generatePgAdminInfo() error {
	printColored("yellow", "Generating pgAdmin connection information...")
	
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
	
	printColored("green", fmt.Sprintf("pgAdmin connection info saved to: %s", infoFile))
	return nil
}

func loadAndValidateConfig() (Config, error) {
	config := Config{
		// HostDB configuration
		HostHost:     os.Getenv("HOST_HOST"),
		HostPort:     os.Getenv("HOST_PORT"),
		HostDB:       os.Getenv("HOST_DB"),
		HostUser:     os.Getenv("HOST_USER"),
		HostPassword: os.Getenv("HOST_PASSWORD"),
		
		// Target PostgreSQL configuration
		TargetHost:     os.Getenv("TARGET_HOST"),
		TargetPort:     os.Getenv("TARGET_PORT"),
		TargetDB:       os.Getenv("TARGET_DB"),
		TargetUser:     os.Getenv("TARGET_USER"),
		TargetPassword: os.Getenv("TARGET_PASSWORD"),
	}
	
	// Set default ports if not provided
	if config.HostPort == "" {
		config.HostPort = "5432"
	}
	if config.TargetPort == "" {
		config.TargetPort = "5432"
	}
	
	// Validate required environment variables
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
		return Config{}, fmt.Errorf("missing required environment variables: %s", strings.Join(missingVars, ", "))
	}
	
	return config, nil
}

func (m *Migrator) Migrate() error {
	printColored("green", "Starting HostDB to PostgreSQL Migration")
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
			printColored("red", fmt.Sprintf("Error in %s: %v", step.name, err))
			return err
		}
	}
	
	fmt.Println()
	printColored("green", "Migration completed successfully!")
	fmt.Println("==================================================")
	fmt.Printf("Backup directory: %s\n", m.backupDir)
	fmt.Println("You can now connect to your PostgreSQL database using pgAdmin")
	fmt.Printf("Connection details are saved in: %s/pgadmin_connection_info.txt\n", m.backupDir)
	
	return nil
}

func main() {
	config, err := loadAndValidateConfig()
	if err != nil {
		printColored("red", err.Error())
		os.Exit(1)
	}
	
	migrator := NewMigrator(config)
	if err := migrator.Migrate(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}