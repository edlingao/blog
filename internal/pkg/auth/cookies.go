package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	AUTHCOOKIE = "Auth"
)

func setCookie(name, value string, c echo.Context) error {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)
	return nil
}

func newCookieWithPath(value, path string) *http.Cookie {
	cookie := &http.Cookie{
		Name:     "sid",
		Value:    value,
		Path:     path,
		MaxAge:   3600,
		Secure:   false,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	}

	return cookie
}

func getCookie(name string, c echo.Context) (string, error) {
	cookie, err := c.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func SetSIdCookie(c echo.Context, sessionToken string) {
	cookie := newCookieWithPath(sessionToken, "/users/")
	c.SetCookie(cookie)
}

func SetAuthCookie(token string, c echo.Context) error {
	return setCookie(AUTHCOOKIE, token, c)
}

func LogoutCookie(c echo.Context) error {
	cookie := &http.Cookie{
		Name:     AUTHCOOKIE,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)
	return nil
}

func IsLoggedIn(c echo.Context) bool {
	value, err := getCookie(AUTHCOOKIE, c)
	if err != nil {
		return false
	}

	valitdToken, err := ValidateToken(value)

	return err == nil && valitdToken.Username != "" && valitdToken.UserID != ""
}

func GetAuthTokenCookie(c echo.Context) (string, error) {
	return getCookie(AUTHCOOKIE, c)
}

func GetAuthTokenHeader(c echo.Context) (string, error) {
	token := c.Request().Header.Get("x-auth")
	if token == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "missing auth token")
	}
	return token, nil
}

func GetAuthToken(c echo.Context) (string, error) {
	// Try to get token from cookie first
	token, err := GetAuthTokenCookie(c)
	if err == nil {
		return token, nil
	}

	// Fallback to header
	return GetAuthTokenHeader(c)
}

func GetAuthClaims(c echo.Context) (*Claims, error) {
	token, err := GetAuthToken(c)
	if err != nil {
		return nil, err
	}
	claims, err := ValidateToken(token)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func GetAuthUserID(c echo.Context) (string, error) {
	claims, err := GetAuthClaims(c)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

func GetAuthUsername(c echo.Context) (string, error) {
	claims, err := GetAuthClaims(c)
	if err != nil {
		return "", err
	}
	return claims.Username, nil
}

func GetSIDCookie(c echo.Context) (string, error) {
	sid, err := getCookie("sid", c)
	if err != nil {
		return "", err
	}

	return sid, nil
}

