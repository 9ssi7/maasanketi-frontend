# Software Engineer Salary Survey Backend

This directory contains the backend service for the Software Engineer Salary Survey platform.

## Status

The backend is currently under initial development. This README will be updated with detailed information as the development progresses.

## Features

- RESTful API for survey data
- Swagger documentation
- Health monitoring
- Google authentication

## API Documentation

The API is documented using Swagger. Once the application is running, you can access the Swagger UI at:

```bash
http://localhost:4004/swagger/
```

To generate or update the Swagger documentation, run:

```bash
make swagger
```

## Authentication

The API uses Google authentication for protected routes. To authenticate:

1. Obtain a Google ID token from the client-side Google Sign-In
2. Include the token in the Authorization header:

   ```txt
   Authorization: Bearer <google-id-token>
   ```

For more details, see the [middleware documentation](middleware/README.md).

## Development

### Prerequisites

- Go 1.24 or higher
- Google OAuth2 client ID (set as `GOOGLE_CLIENT_ID` environment variable)

### Running the Application

```bash
make up
```

The API will be available at `http://localhost:4004`.

Check back for updates as development continues.

## Survey Graph Generation Cron Job

The application includes a worker service that automatically generates graphs from survey responses when surveys are completed.

### How It Works

1. The worker runs as a separate service that checks for completed surveys every hour.
2. For each completed survey, it generates graph data based on the survey responses.
3. The graphs are generated according to their configuration (bar, line, pie, etc.) and value strategy (count, sum, avg).

### Configuration

The worker uses the following components:

- `GraphGenerationService`: The core service that manages the graph generation workflow.
- `gocron`: The scheduling library that manages the periodic execution of the job.

### Running the Worker

#### Using Docker Compose

The worker is included in the `docker-compose.yml` file and will start automatically when you run:

```
docker-compose up
```

#### Running Separately

To run the worker separately:

```
go run cmd/worker/main.go
```

### Customizing the Schedule

By default, the worker checks for completed surveys every hour. To change this schedule, modify the scheduling configuration in `cmd/worker/main.go`:

```go
// Example: Run every 30 minutes
scheduler.Every(30).Minutes().Do(...)

// Example: Run every day at midnight
scheduler.Every(1).Day().At("00:00").Do(...)
```

### Monitoring

The worker logs its activities using the application's logging system. Monitor these logs to ensure the graph generation process is working correctly.
