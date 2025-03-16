# Middleware

This directory contains middleware components for the Software Engineer Salary Survey API.

## Google Authentication Middleware

The Google Authentication middleware verifies Google ID tokens for protected routes.

### Configuration

The middleware requires a Google OAuth2 client ID, which can be set via the `GOOGLE_CLIENT_ID` environment variable.

To create a Google OAuth2 client ID:

1. Go to the [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Navigate to "APIs & Services" > "Credentials"
4. Click "Create Credentials" > "OAuth client ID"
5. Configure the OAuth consent screen
6. Create a Web application client ID
7. Note the client ID for use in the application

### Usage

To use the middleware in your application:

```go
// Apply to a group of routes
userRoutes := app.Group("/api/v1/user")
userRoutes.Use(middleware.GoogleAuth())
userRoutes.Get("/profile", handler.GetUserProfile)

// Or with custom configuration
userRoutes.Use(middleware.GoogleAuth(middleware.GoogleAuthConfig{
    ClientID:  "your-client-id.apps.googleusercontent.com",
    SkipPaths: []string{"/public"},
}))
```

### Client-Side Integration

On the client side, you need to:

1. Implement Google Sign-In
2. Get the ID token from the Google Sign-In response
3. Include the token in the Authorization header for API requests:

```javascript
// Example using fetch API
fetch('http://localhost:4004/api/v1/user/profile', {
  headers: {
    'Authorization': `Bearer ${googleIdToken}`
  }
})
.then(response => response.json())
.then(data => console.log(data));
```

### Accessing User Information in Handlers

The middleware sets the following user information in the Fiber context:

- `user`: The full ID token payload
- `user_id`: The user's Google ID (sub claim)
- `email`: The user's email address
- `name`: The user's name
- `picture`: The URL of the user's profile picture

You can access this information in your handlers:

```go
func GetUserProfile(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(string)
    email := c.Locals("email").(string)
    
    // Use the user information
    return c.JSON(fiber.Map{
        "id": userID,
        "email": email,
    })
}
```
