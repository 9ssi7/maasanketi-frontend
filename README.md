# Software Engineer Salary Survey

A monorepo containing the frontend and backend services for the Software Engineer Salary Survey platform. This project aims to improve salary transparency in the tech industry.

![image](modules/web/public/image.png)

## Overview

This repository contains all the code needed to run the Software Engineer Salary Survey platform:

- **Frontend**: React application for creating, participating in, and viewing results of salary surveys
- **Backend**: Go-based API server for data storage and processing

## Repository Structure

```txt
maasanketi/
├── modules/
│   ├── web/         # Frontend React application
│   └── backend/     # Backend Go application
└── README.md        # This file
```

## Getting Started

Each module has its own README with specific setup instructions:

- [Frontend Documentation](modules/web/README.md)
- [Backend Documentation](modules/backend/README.md)

## Development

For local development, you'll need to set up both the frontend and backend services. See the respective README files for detailed instructions.

### Installing golang-migrate

You can install the golang-migrate CLI tool using the provided Makefile command:

```bash
make install-migrate
```

This will install the latest version of golang-migrate with PostgreSQL support. Make sure your Go bin directory is in your PATH.

### Setting Up the Development Environment

1. Start the PostgreSQL database:

   ```bash
   make db-up
   ```

2. Run database migrations:

   ```bash
   make migrate-up
   ```

   Or use the all-in-one setup command:

   ```bash
   make dev-setup
   ```

### Running the Application

```bash
make up
```

The API will be available at `http://localhost:4004`.

### Database Migrations

The project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations. Migration files are stored in the `resources/migrations` directory.

#### Creating a New Migration

```bash
make migrate-create NAME=create_new_table
```

This will create two new files in the `resources/migrations` directory:

- `NNNNNN_create_new_table.up.sql`: Contains SQL statements to apply the migration
- `NNNNNN_create_new_table.down.sql`: Contains SQL statements to revert the migration

#### Applying Migrations

```bash
make migrate-up
```

#### Rolling Back Migrations

```bash
make migrate-down
```

#### Other Migration Commands

Run `make help` to see all available migration commands.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
