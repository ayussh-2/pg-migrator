package main

import (
	"log"
	"os"
	"pg-migrator/migrator"
	"pg-migrator/utils"
)

func main() {
	config, err := utils.LoadConfig()
	if err != nil {
		migrator.PrintColored("red", err.Error())
		os.Exit(1)
	}

	m := migrator.NewMigrator(config)
	if err := m.Migrate(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}
