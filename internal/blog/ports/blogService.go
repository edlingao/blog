package ports

import "github.com/labstack/echo/v4"

type BlogService interface {
	IndexHandler(echo.Context) error
	AboutHandler(echo.Context) error
	PostDetailsHandler(echo.Context) error
	PostEventsHandler(echo.Context) error
	PortfolioHandler(echo.Context) error
	PostCommentsHandler(echo.Context) error
	AdminIndexHandler(echo.Context) error
	AdminPostsHandler(echo.Context) error
	AdminCommentsHandler(echo.Context) error
	AdminDeleteCommentHandler(echo.Context) error
	AdminDeletePostHandler(echo.Context) error
	AdminToggleCommentsHandler(echo.Context) error
}
