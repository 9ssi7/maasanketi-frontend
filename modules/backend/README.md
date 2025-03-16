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
