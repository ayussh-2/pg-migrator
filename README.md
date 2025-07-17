# PostgreSQL Migration Tool

This Go application facilitates the migration of a PostgreSQL database from a host instance to a target PostgreSQL instance. It automates the process of dumping the source database and restoring it to the destination, providing checks and verifications along the way.

---

## What it Does

This tool performs the following key operations:

1. **PostgreSQL Tool Check:** Verifies the presence of essential PostgreSQL client tools (`pg_dump` and `psql`) on your system.
2. **Database Connection Testing:** Ensures successful connections to both the source (HostDB) and target PostgreSQL databases.
3. **Database Dump Creation:** Generates a SQL dump of the source HostDB database using `pg_dump`. This dump is stored in a time-stamped directory for easy organization and potential recovery.
4. **Database Restore:** Restores the created SQL dump to the target PostgreSQL database using `psql`.
5. **Migration Verification:** Compares the number of tables in the source and target databases to provide a basic verification of the migration's success.
6. **pgAdmin Connection Information Generation:** Creates a text file containing connection details for the newly migrated target database, making it easy to connect via pgAdmin or other clients.

---

## What is Used

The application is written in Go and leverages the following:

-   `database/sql`: Go's standard library for interacting with SQL databases.
-   `github.com/lib/pq`: A pure Go PostgreSQL driver for `database/sql`.
-   `os`: For interacting with the operating system, such as setting environment variables and creating directories.
-   `os/exec`: To execute external commands, specifically `pg_dump` and `psql`.
-   `path/filepath`: For manipulating file paths.
-   `fmt`: For formatted I/O.
-   `log`: For logging fatal errors.
-   `time`: For generating timestamps.
-   `strings`: For string manipulation.

---

## How to Install

To use this tool, you need to have Go and PostgreSQL client tools installed on your system.

### Prerequisites

-   **Go:** Ensure you have Go installed (version 1.24 or higher is recommended). You can download it from the [official Go website](https://golang.org/dl/).
-   **PostgreSQL Client Tools:** You need `pg_dump` and `psql` available in your system's PATH. These are typically included when you install PostgreSQL or its client packages.

    For example, on Debian/Ubuntu:

    ```bash
    sudo apt-get install postgresql-client
    ```

---

### Running the Application

1. Clone the repository (or save the provided code as a `.go` file, e.g., `main.go`).
2. Navigate to the project directory in your terminal.
3. Download dependencies:

    ```bash
    go mod init pg-migrator
    go mod tidy
    ```

4. Run the script:

    ```bash
    go run main.go
    ```

---

## How to Use

The application requires database connection details to be provided via environment variables.

### Configuration

Copy the `.sample.env` file to `.env` (ensure it is saved in UTF-8 format).

#### Source (HostDB) Configuration

-   `HOST_HOST`: Hostname or IP address of the source PostgreSQL database.
-   `HOST_PORT`: Port of the source PostgreSQL database (defaults to 5432 if not set).
-   `HOST_DB`: Database name on the source.
-   `HOST_USER`: Username for connecting to the source database.
-   `HOST_PASSWORD`: Password for the source database user.

#### Target PostgreSQL Configuration

-   `TARGET_HOST`: Hostname or IP address of the target PostgreSQL database.
-   `TARGET_PORT`: Port of the target PostgreSQL database (defaults to 5432 if not set).
-   `TARGET_DB`: Database name on the target where data will be restored.
-   `TARGET_USER`: Username for connecting to the target database.
-   `TARGET_PASSWORD`: Password for the target database user.
