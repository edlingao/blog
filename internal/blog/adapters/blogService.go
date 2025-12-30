package adapters

import (
	"errors"
	"log"
	"time"

	"github.com/edlingao/internal/blog/core"
	"github.com/edlingao/internal/blog/ports"
	"github.com/edlingao/internal/pkg/auth"
	"github.com/edlingao/internal/pkg/web"
	"github.com/edlingao/web/template/admin"
	"github.com/edlingao/web/template/blogs"
	"github.com/edlingao/web/template/pages"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type BlogService struct {
	root                 *echo.Group
	blogRepo             ports.BlogRepository
	commentRepo          ports.CommentsRepository
	commentEventsService ports.CommentsChannelRepo
}

func NewBlogService(
	root *echo.Group,
	blogRepo ports.BlogRepository,
	commentRepo ports.CommentsRepository,
	commentEventsService ports.CommentsChannelRepo,
) *BlogService {
	service := &BlogService{
		root:                 root,
		blogRepo:             blogRepo,
		commentRepo:          commentRepo,
		commentEventsService: commentEventsService,
	}

	service.root.GET("/", service.IndexHandler)
	service.root.GET("/post/:title", service.PostDetailsHandler)
	service.root.GET("/post/:title/stream", service.PostEventsHandler)
	service.root.GET("/about", service.AboutHandler)
	service.root.GET("/portfolio", service.PortfolioHandler)

	service.root.POST("/post/:title/comments", service.PostCommentsHandler)

	protected := service.root.Group("/admin", auth.AuthMiddleware)
	protected.GET("/posts", service.AdminPostsHandler)
	protected.GET("/comments", service.AdminCommentsHandler)
	protected.DELETE("/comment/:id", service.AdminDeleteCommentHandler)
	protected.DELETE("/post/:id", service.AdminDeletePostHandler)
	protected.POST("/post/:id/comments/toggle", service.AdminToggleCommentsHandler)

	go service.commentEventsService.Start()

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
		return web.Render(c, blogs.Detail(blogs.BlogDetailProps{
			Error: errors.New("Post not found"),
		}), 404)
	}

	tags, err := service.blogRepo.GetTagsByBlogID(post.ID)
	if err != nil {
		return web.Render(c, blogs.Detail(blogs.BlogDetailProps{
			Error: errors.New("Error: " + err.Error()),
		}), 404)
	}
	post.Tags = tags

	comments, err := service.commentRepo.GetAllCommentsByPostIDWithChildren(post.ID)
	if err != nil {
		log.Print(err)
		return web.Render(c, blogs.Detail(blogs.BlogDetailProps{
			Error: errors.New("Error: " + err.Error()),
		}), 404)
	}
	post.Comments = comments

	return web.Render(c, blogs.Detail(blogs.BlogDetailProps{
		Post:            post,
		HTMLContent:     post.GetContent(),
		CommentsEnabled: post.CommentsAvailable,
		Comments:        post.Comments,
	}), 200)
}

func (service *BlogService) PortfolioHandler(c echo.Context) error {
	return web.Render(c, pages.Portfolio(), 200)
}

func (service *BlogService) PostEventsHandler(c echo.Context) error {
	clientID := uuid.NewString()
	title := c.QueryParam("title")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache, no-transform")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	client := core.NewClient(clientID, title, w)
	go service.commentEventsService.AddSubscriber(client)
	for {
		select {
		case <-c.Request().Context().Done():
			go service.commentEventsService.UnSubscribe(client)
			c.Request().Body.Close()
			return nil
		case <-ticker.C:
			go client.SendEvent(&core.Event{
				Title: title,
				Data:  map[string]any{"message": "ping"},
			})
		}
	}
}

func (service *BlogService) PostCommentsHandler(c echo.Context) error {
	postTitle := c.Param("title")
	author := c.FormValue("author")
	content := c.FormValue("content")
	replyTo := c.FormValue("reply_to")

	post, err := service.blogRepo.GetByTitle(postTitle)
	if err != nil {
		return web.Render(c, blogs.Detail(blogs.BlogDetailProps{
			Error: errors.New("Post not found"),
		}), 404)
	}

	if !post.CommentsAvailable {
		c.Response().Header().Add("HX-Refresh", "true")
		return web.Render(c, blogs.Detail(blogs.BlogDetailProps{
			Post:  post,
			Error: errors.New("Comments are disabled for this post"),
		}), 500)
	}

	var comment *core.Comment
	comment = core.NewComment(post.ID, author, content)
	if replyTo != "" {
		parentComment, err := service.commentRepo.GetCommentByID(replyTo)
		if err != nil {
			log.Print(err)
			return web.Render(c, blogs.Detail(blogs.BlogDetailProps{
				Post:  post,
				Error: errors.New("Parent comment not found: " + err.Error()),
			}), 500)
		}
		comment.AddParentComment(parentComment.ID)
		comment, err = service.commentRepo.ReplyComment(parentComment.ID, comment)
	} else {
		comment, err = service.commentRepo.SaveComment(comment)
	}

	if err != nil {
		log.Print(err)
		return web.Render(c, blogs.Detail(blogs.BlogDetailProps{
			Post:  post,
			Error: errors.New("Failed to save comment: " + err.Error()),
		}), 500)
	}

	event := core.NewEvent("new_comment", postTitle, map[string]any{
		"author":  author,
		"content": content,
		"id":      comment.ID,
		"parent":  comment.CommentID,
	})

	service.commentEventsService.BroadcastEvent(event)

	return c.String(200, "Comment added successfully")
}

func (service *BlogService) AdminPostsHandler(c echo.Context) error {
	posts, err := service.blogRepo.GetAllBlogs()
	if err != nil {
		return web.Render(c, admin.Posts(admin.AdminPostsProps{
			Message: "Failed to load posts: " + err.Error(),
		}), 500)
	}

	return web.Render(c, admin.Posts(admin.AdminPostsProps{
		Posts: posts,
	}), 200)
}

func (service *BlogService) AdminCommentsHandler(c echo.Context) error {
	posts, err := service.blogRepo.GetAllBlogs()
	if err != nil {
		return web.Render(c, admin.Comments(admin.AdminCommentsProps{
			Message: "Failed to load posts: " + err.Error(),
		}), 500)
	}

	var allCommentsByPost []admin.CommentsByPost
	for _, post := range posts {
		comments, err := service.commentRepo.GetAllCommentsByPostIDWithChildren(post.ID)
		if err != nil {
			return web.Render(c, admin.Comments(admin.AdminCommentsProps{
				Message: "Failed to load comments for post " + post.Title + ": " + err.Error(),
			}), 500)
		}
		allCommentsByPost = append(allCommentsByPost, admin.CommentsByPost{
			PostID:    post.ID,
			PostTitle: post.Title,
			Comments:  comments,
		})

	}

	return web.Render(c, admin.Comments(admin.AdminCommentsProps{
		CommentsByPost: allCommentsByPost,
	}), 200)
}

func (service *BlogService) AdminIndexHandler(c echo.Context) error {
	posts, err := service.blogRepo.GetAllBlogs()
	if err != nil {
		return web.Render(c, admin.Posts(admin.AdminPostsProps{
			Message: "Failed to load posts: " + err.Error(),
		}), 500)
	}

	return web.Render(c, admin.Posts(admin.AdminPostsProps{
		Posts: posts,
	}), 200)
}

func (service *BlogService) AdminDeleteCommentHandler(c echo.Context) error {
	commentID := c.Param("id")
	err := service.commentRepo.DeleteComment(commentID)
	if err != nil {
		return web.Render(c, admin.Comments(admin.AdminCommentsProps{
			Message: "Failed to delete comment: " + err.Error(),
		}), 500)
	}

	return web.Render(c, admin.Comments(admin.AdminCommentsProps{
		Message: "Comment deleted successfully",
	}), 200)
}

func (service *BlogService) AdminDeletePostHandler(c echo.Context) error {
	commentID := c.Param("id")
	err := service.blogRepo.DeleteBlog(commentID)
	if err != nil {
		return web.Render(c, admin.Posts(admin.AdminPostsProps{
			Message: "Failed to delete post: " + err.Error(),
		}), 500)
	}

	return web.Render(c, admin.Comments(admin.AdminCommentsProps{
		Message: "Comment deleted successfully",
	}), 200)
}

func (service *BlogService) AdminToggleCommentsHandler(c echo.Context) error {
	postID := c.Param("id")
	post, err := service.blogRepo.GetBlogByID(postID)
	if err != nil {
		return web.Render(c, admin.Posts(admin.AdminPostsProps{
			Message: "Failed to load post: " + err.Error(),
		}), 500)
	}

	post.ToggleCommentsAvailability()
	post, err = service.blogRepo.Update(post)
	if err != nil {
		return web.Render(c, admin.Posts(admin.AdminPostsProps{
			Message: "Failed to update post: " + err.Error(),
		}), 500)
	}

	posts, err := service.blogRepo.GetAllBlogs()
	if err != nil {
		return web.Render(c, admin.Posts(admin.AdminPostsProps{
			Message: "Failed to load posts: " + err.Error(),
		}), 500)
	}

	return web.Render(c, admin.Posts(admin.AdminPostsProps{
		Message: "Post updated successfully",
		Posts:   posts,
	}), 200)
}
