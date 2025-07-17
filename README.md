# PostgreSQL Migration Tool

This Go application facilitates the migration of a PostgreSQL database from a host instance to a target PostgreSQL instance. It supports both traditional host-based migrations and Docker-based deployments, automatically detecting the best approach based on your configuration.

---

## What it Does

This tool performs the following key operations:

1. **PostgreSQL Tool Check:** Verifies the presence of essential PostgreSQL client tools (`pg_dump` and `psql`) either on your host system or within Docker containers.
2. **Database Connection Testing:** Ensures successful connections to both the source (host) and target PostgreSQL databases.
3. **Database Dump Creation:** Generates a SQL dump of the source database using `pg_dump`. This dump is stored in a time-stamped directory for easy organization and potential recovery.
4. **Database Restore:** Restores the created SQL dump to the target PostgreSQL database using `psql`.
5. **Migration Verification:** Compares the number of tables in the source and target databases to provide a basic verification of the migration's success.
6. **pgAdmin Connection Information Generation:** Creates a text file containing connection details for both source and target databases, making it easy to connect via pgAdmin or other clients.

---

## Features

-   **Hybrid Deployment Support:** Works with any combination of local, remote, and Docker-based PostgreSQL instances
-   **External Database Support:** Seamlessly connects to cloud databases like Neon, AWS RDS, etc.
-   **Docker Integration:** Automatically uses Docker containers for PostgreSQL tools when available
-   **Smart Tool Detection:** Falls back gracefully between Docker and host-based tools
-   **Comprehensive Logging:** Color-coded output for easy monitoring of migration progress
-   **Backup Management:** Creates timestamped backup directories for easy organization

---

## What is Used

The application is written in Go and leverages the following:

-   `database/sql`: Go's standard library for interacting with SQL databases
-   `github.com/lib/pq`: A pure Go PostgreSQL driver for `database/sql`
-   `github.com/joho/godotenv`: For loading environment variables from `.env` files
-   `os/exec`: To execute external commands (`pg_dump`, `psql`) and Docker commands
-   Standard Go libraries: `os`, `path/filepath`, `fmt`, `log`, `time`, `strings`

---

## How to Install

### Prerequisites

**Required:**

-   **Go:** Version 1.19 or higher. Download from the [official Go website](https://golang.org/dl/)
-   **Docker:** If using Docker-based PostgreSQL instances. Download from [Docker website](https://www.docker.com/get-started)

**Optional (choose based on your setup):**

-   **PostgreSQL Client Tools:** Required if running tools on host machine
    -   Windows: Download from [PostgreSQL Downloads](https://www.postgresql.org/download/windows/)
    -   Linux: `sudo apt-get install postgresql-client` (Ubuntu/Debian) or equivalent
    -   macOS: `brew install postgresql`

### Installation Steps

1. Clone the repository:

    ```bash
    git clone <your-forked-repository-url>
    cd pg-migrator
    ```

2. Initialize Go modules:

    ```bash
    go mod init pg-migrator
    go mod tidy
    ```

3. Build the application (optional):
    ```bash
    go build
    ```

---

## How to Use

### Configuration

Create a `.env` file in the project root with your database connection details:

```env
# Source Database Configuration
HOST_HOST=your-source-host                    # e.g., neon-host.com or localhost
HOST_PORT=5432                                # Default: 5432
HOST_DB=your-source-database
HOST_USER=your-source-username
HOST_PASSWORD=your-source-password
HOST_SSLMODE=require                          # For cloud DBs, use 'disable' for local
HOST_DOCKER_CONTAINER=                        # Optional: Docker container name if source is in Docker

# Target Database Configuration
TARGET_HOST=localhost                         # Or your target host
TARGET_PORT=5432                             # Default: 5432
TARGET_DB=your-target-database
TARGET_USER=your-target-username
TARGET_PASSWORD=your-target-password
TARGET_SSLMODE=disable                       # Usually 'disable' for local Docker
TARGET_DOCKER_CONTAINER=your-postgres-container  # Docker container name

# Tool Execution Strategy (Optional)
USE_DOCKER_FOR_TOOLS=true                    # Force using Docker for tools
```

### Configuration Examples

#### Example 1: Neon (Cloud) → Local Docker

```env
# Source: Neon Database (External)
HOST_HOST=ep-spring-bush-a1if0ihi-pooler.ap-southeast-1.aws.neon.tech
HOST_DB=your_database
HOST_USER=your_user
HOST_PASSWORD=your_password
HOST_SSLMODE=require

# Target: Local Docker PostgreSQL
TARGET_HOST=localhost
TARGET_DB=your_database
TARGET_USER=postgres
TARGET_PASSWORD=postgres
TARGET_SSLMODE=disable
TARGET_DOCKER_CONTAINER=my-postgres-container
```

#### Example 2: Docker → Docker

```env
# Source: Docker PostgreSQL
HOST_HOST=localhost
HOST_DB=source_db
HOST_USER=postgres
HOST_PASSWORD=postgres
HOST_DOCKER_CONTAINER=source-postgres

# Target: Docker PostgreSQL
TARGET_HOST=localhost
TARGET_DB=target_db
TARGET_USER=postgres
TARGET_PASSWORD=postgres
TARGET_DOCKER_CONTAINER=target-postgres
```

#### Example 3: Local → Remote

```env
# Source: Local PostgreSQL
HOST_HOST=localhost
HOST_DB=local_db
HOST_USER=postgres
HOST_PASSWORD=postgres
HOST_SSLMODE=disable

# Target: Remote PostgreSQL (e.g., AWS RDS)
TARGET_HOST=your-rds-endpoint.amazonaws.com
TARGET_DB=production_db
TARGET_USER=your_user
TARGET_PASSWORD=your_password
TARGET_SSLMODE=require
```

### Running the Migration

Execute the migration with:

```bash
# Using go run
go run main.go

# Or using built binary
./pg-migrator
```

### Expected Output

```
Starting HostDB to PostgreSQL Migration
==================================================
Checking for Docker...
Checking pg_dump on host machine...
pg_dump found on host machine.
Checking psql in target Docker container...
psql found in target container.
Testing HostDB connection...
HostDB connection successful.
Testing Target PostgreSQL connection...
Target PostgreSQL connection successful.
Creating database dump from host machine...
Database dump created successfully: host_migration_20250717_164728\host_dump.sql
Dump file size: 2.45 MB
Restoring dump to target via Docker container...
Database restore completed successfully.
Verifying migration...
HostDB tables: 15
Target PostgreSQL tables: 15
Table count verification passed.
Generating pgAdmin connection information...
Connection info saved to: host_migration_20250717_164728\pgadmin_connection_info.txt

Migration completed successfully!
==================================================
Backup directory: host_migration_20250717_164728
You can now connect to your PostgreSQL database using pgAdmin
Connection details are saved in: host_migration_20250717_164728/pgadmin_connection_info.txt
```

---

## Docker Support

The tool automatically detects and adapts to your Docker setup:

-   **Auto-detection:** Automatically uses Docker containers when specified
-   **Network Access:** Docker containers can access external databases (like Neon, AWS RDS)
-   **Tool Flexibility:** Uses Docker-based `pg_dump`/`psql` when containers are available
-   **Fallback Support:** Falls back to host tools if Docker tools aren't available

### Docker Container Requirements

Your PostgreSQL Docker containers should:

1. Be running and accessible
2. Have PostgreSQL client tools installed (`pg_dump`, `psql`)
3. Be able to access external networks (default Docker behavior)

---

## Troubleshooting

### Common Issues

1. **"pg_dump not found"**

    - Install PostgreSQL client tools on your host, or
    - Ensure your Docker containers have PostgreSQL tools installed

2. **"Connection failed"**

    - Verify database credentials in `.env`
    - Check network connectivity to remote databases
    - Ensure Docker containers are running

3. **"Docker command not found"**

    - Install Docker if using Docker-based configurations
    - Or remove Docker container names from `.env` to use host tools

4. **"Password authentication failed"**
    - Double-check credentials in `.env`
    - Ensure special characters in passwords are properly escaped

### SSL Modes

-   **require:** For secure connections (cloud databases)
-   **disable:** For local development databases
-   **prefer:** Attempts SSL, falls back to non-SSL

---

## Project Structure

```
pg-migrator/
├── main.go                 # Application entry point
├── utils/
│   └── env.go             # Configuration loading
├── migrator/
│   ├── migrator.go        # Main migration orchestrator
│   ├── tools.go           # Tool availability checking
│   ├── connection.go      # Database connection testing
│   ├── dump.go            # Database dump and restore
│   └── verify.go          # Migration verification
├── .env                   # Your configuration (create this)
├── .sample.env           # Sample configuration
└── README.md             # This file
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details on how to get started.

---

## Support

If you encounter any issues or have questions, please:

1. Check the [Troubleshooting](#troubleshooting) section
2. Search existing [GitHub Issues](../../issues)
3. Create a new issue if your problem isn't covered

---

## Acknowledgments

-   Thanks to the PostgreSQL community for excellent tools
-   Built with Go and love ❤️

---

⭐ **If this project helped you, please consider giving it a star!** ⭐
