package migrator

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// testConnection tests a single database connection.
func (m *Migrator) testConnection(host, port, dbname, user, password, sslmode, name string) error {
	PrintColored("yellow", fmt.Sprintf("Testing %s connection...", name))

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	fmt.Printf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s\n",
	host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %v", name, err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping %s: %v", name, err)
	}

	PrintColored("green", fmt.Sprintf("%s connection successful.", name))
	return nil
}

// testConnections tests both source and target database connections.
func (m *Migrator) testConnections() error {
	if err := m.testConnection(
		m.config.HostHost, m.config.HostPort, m.config.HostDB,
		m.config.HostUser, m.config.HostPassword, m.config.HostSSLMode, "HostDB",
	); err != nil {
		return err
	}

	if err := m.testConnection(
		m.config.TargetHost, m.config.TargetPort, m.config.TargetDB,
		m.config.TargetUser, m.config.TargetPassword, m.config.TargetSSLMode, "Target PostgreSQL",
	); err != nil {
		return err
	}

	return nil
}

// getTableCount gets the number of tables in a database.
func (m *Migrator) getTableCount(host, port, dbname, user, password, sslmode string) (int, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

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
