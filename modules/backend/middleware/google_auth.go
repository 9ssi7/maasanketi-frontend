package middleware

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	log "github.com/mstrYoda/maasanketi.co/pkg/logger"
	"google.golang.org/api/idtoken"
)

type GoogleAuthConfig struct {
	// ClientID is the Google OAuth2 client ID
	ClientID string
	// SkipPaths are the paths that should skip authentication
	SkipPaths []string
}

func DefaultGoogleAuthConfig() GoogleAuthConfig {
	return GoogleAuthConfig{
		ClientID:  os.Getenv("GOOGLE_CLIENT_ID"),
		SkipPaths: []string{"/healthcheck", "/metrics", "/swagger", "/monitor"},
	}
}

func GoogleAuth(config ...GoogleAuthConfig) fiber.Handler {
	cfg := DefaultGoogleAuthConfig()

	if len(config) > 0 {
		cfg = config[0]
	}

	if cfg.ClientID == "" {
		log.Logger().Warn("Google client ID is not set, authentication will fail")
	}

	return func(c *fiber.Ctx) error {
		path := c.Path()
		for _, skipPath := range cfg.SkipPaths {
			if strings.HasPrefix(path, skipPath) {
				return c.Next()
			}
		}

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}
		idToken := parts[1]
		payload, err := verifyGoogleIDToken(idToken, cfg.ClientID)
		if err != nil {
			log.Logger().Error(fmt.Sprintf("Failed to verify Google ID token: %v", err))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		c.Locals("user", payload)
		c.Locals("email", payload.Claims["email"])
		c.Locals("user_id", payload.Claims["sub"])
		c.Locals("name", payload.Claims["name"])
		c.Locals("picture", payload.Claims["picture"])
		return c.Next()
	}
}

func verifyGoogleIDToken(idToken, clientID string) (*idtoken.Payload, error) {
	ctx := context.Background()
	validator, err := idtoken.NewValidator(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator: %v", err)
	}

	payload, err := validator.Validate(ctx, idToken, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}

	return payload, nil
}
