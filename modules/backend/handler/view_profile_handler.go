package handler

import (
	"github.com/gofiber/fiber/v2"
)

// UserProfile represents the user profile data
type UserProfile struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// @Summary View user profile
// @Description View the profile of the authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserProfile
// @Failure 401 {object} map[string]string
// @Router /api/v1/user/profile [get]
func ViewProfile(c *fiber.Ctx) error {
	// Get user information from context (set by the GoogleAuth middleware)
	userID := c.Locals("user_id")
	email := c.Locals("email")
	name := c.Locals("name")
	picture := c.Locals("picture")

	// Create user profile
	profile := UserProfile{
		ID:      userID.(string),
		Email:   email.(string),
		Name:    name.(string),
		Picture: picture.(string),
	}

	return c.JSON(profile)
}
