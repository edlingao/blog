package adapters

import (
	"github.com/edlingao/internal/blog/core"
	"github.com/edlingao/internal/blog/ports"
	"github.com/edlingao/internal/pkg/web"
	"github.com/edlingao/web/template/blogs"
	"github.com/edlingao/web/template/pages"
	"github.com/labstack/echo/v4"
)

type BlogService struct {
	root     *echo.Group
	blogRepo ports.BlogRepository
}

func NewBlogService(
	root *echo.Group,
	blogRepo ports.BlogRepository,
) *BlogService {
	service := &BlogService{
		root:     root,
		blogRepo: blogRepo,
	}

	service.root.GET("/", service.IndexHandler)
	service.root.GET("/post/:title", service.PostDetailsHandler)
	service.root.GET("/about", service.AboutHandler)
	service.root.GET("/portfolio", service.PortfolioHandler)

	return service
}

func (service *BlogService) AboutHandler(c echo.Context) error {
	return web.Render(c, pages.About(), 200)
}

func (service *BlogService) IndexHandler(c echo.Context) error {
	tag := c.QueryParam("tag")
	tags, err := service.blogRepo.GetTagsWithBlogCount()
	if err != nil {
		return web.Render(c, blogs.Index(blogs.PostsIndexProps{}), 200)
	}

	var posts []core.Blog
	if tag != "" {
		posts, err = service.blogRepo.GetAllBlogsByTag(tag)
		if err != nil {
			return web.Render(c, blogs.Index(blogs.PostsIndexProps{}), 200)
		}
	} else {
		posts, err = service.blogRepo.GetAllBlogs()
		if err != nil {
			return web.Render(c, blogs.Index(blogs.PostsIndexProps{}), 200)
		}
	}

	return web.Render(c, blogs.Index(blogs.PostsIndexProps{
		Tags:      tags,
		Posts:     posts,
		ActiveTag: tag,
	}), 200)
}

func (service *BlogService) PostDetailsHandler(c echo.Context) error {
	postTitle := c.Param("title")
	post, err := service.blogRepo.GetByTitle(postTitle)
	if err != nil {
		return c.NoContent(404)
	}

	tags, err := service.blogRepo.GetTagsByBlogID(post.ID)
	if err != nil {
		return c.NoContent(500)
	}
	post.Tags = tags

	return web.Render(c, blogs.Detail(blogs.BlogDetailProps{
		Post:            post,
		HTMLContent:     post.GetContent(),
		CommentsEnabled: true,
	}), 200)
}

func (service *BlogService) PortfolioHandler(c echo.Context) error {
	return web.Render(c, pages.Portfolio(), 200)
}
