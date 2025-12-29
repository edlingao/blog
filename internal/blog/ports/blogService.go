package ports

import "github.com/labstack/echo/v4"

type BlogService interface {
	IndexHandler(echo.Context) error
	AboutHandler(echo.Context) error
	PostDetailsHandler(echo.Context) error
	PortfolioHandler(echo.Context) error
}
