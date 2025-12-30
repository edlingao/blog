package auth

import (
	"github.com/edlingao/internal/pkg/database"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		value, err := GetAuthToken(c)
		if err != nil {
			c.Response().Header().Add("HX-Redirect", "/login")
			return c.Redirect(302, "/login")
		}

		claims, err := ValidateToken(value)
		if err != nil {
			c.Response().Header().Add("HX-Redirect", "/login")
			return c.Redirect(302, "/login")
		}

		db := database.New()
		defer db.Close()

		userExists, err := checkUserExists(db, claims.UserID)
		if err != nil || !userExists {
			c.Response().Header().Add("HX-Redirect", "/login")
			return c.Redirect(302, "/login")
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		return next(c)
	}
}

func APIAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("x-auth")
		db := database.New()
		defer db.Close()

		if token == "" {
			return c.JSON(401, map[string]string{"error": "missing auth token"})
		}

		claims, err := ValidateToken(token)
		if err != nil {
			return c.JSON(401, map[string]string{"error": "invalid token"})
		}

		userExists, err := checkUserExists(db, claims.UserID)
		if err != nil || !userExists {
			return c.JSON(401, map[string]string{"error": "user does not exist"})
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		return next(c)
	}
}

func checkUserExists(db *sqlx.DB, userID string) (bool, error) {
	var exists bool
	err := db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", userID)
	if err != nil {
		return false, err
	}
	return exists, nil
}
