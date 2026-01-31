package middleware

import (
	"cqs-kanban/internal/service"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func JWTMiddleware(authService service.AuthService, whiteList []string) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Skip middleware for whitelisted routes
		for _, route := range whiteList {
			if c.Path() == route {
				return c.Next()
			}
		}

		// Get the JWT token from the request
		authHeader := c.Get("Authorization")

		// Check if auth header exists
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Authorization required",
				"message": "Missing Authorization header",
			})
		}

		// Bypass if auth token is super admin
		if authHeader == "Basic 17c4520f6cfd1ab53d8745e84681eb49" {
			c.Locals("user_id", 0)
			c.Locals("username", "super_admin")
			c.Locals("is_admin", true)
			c.Locals("department_id")
			fmt.Println("Super Admin Access Granted")
			return c.Next()
		}

		// Check if auth header format is valid
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Invalid authorization format",
				"message": "Authorization header must be in format: Bearer {token}",
			})
		}

		// Validate token
		tokenString := parts[1]
		claims, err := authService.ValidateToken(c.RequestCtx(), tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Invalid token",
				"message": err.Error(),
			})
		}

		if claims.UserID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Invalid user",
				"message": "User not found",
			})
		}

		// Set user info in context
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("department_id", claims.DepartmentID)
		fmt.Println("Authenticated user ID:", claims.UserID, "Username:", claims.Username)

		// Continue to next handler
		return c.Next()
	}
}
