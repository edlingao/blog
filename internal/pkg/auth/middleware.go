package auth

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		value, err := GetAuthToken(c)
		if err != nil {
			fmt.Println("DEBUG: No cookie found:", err)
			return c.Redirect(302, "/users/login")
		}
		fmt.Println("DEBUG: Cookie value length:", len(value))

		claims, err := ValidateToken(value)
		if err != nil {
			fmt.Println("DEBUG: Token validation failed:", err)
			return c.Redirect(302, "/users/login")
		}
		fmt.Println("DEBUG: Token valid, user:", claims.Username)

		return next(c)
	}
}

func APIAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("x-auth")
		if token == "" {
			return c.JSON(401, map[string]string{"error": "missing auth token"})
		}

		claims, err := ValidateToken(token)
		if err != nil {
			return c.JSON(401, map[string]string{"error": "invalid token"})
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		return next(c)
	}
}

