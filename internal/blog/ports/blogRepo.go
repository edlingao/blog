package ports

import "github.com/edlingao/internal/blog/core"

type BlogRepository interface {
	Save(blog *core.Blog) (*core.Blog, error)
	Update(blog *core.Blog) (*core.Blog, error)
	GetByTitle(postTitle string) (*core.Blog, error)
	AddTagsToBlog(blogID string, tags []string) error
	RemoveTagFromBlog(blogID string, tagID string) error
	GetTagsByBlogID(blogID string) ([]core.Tag, error)
	GetTagsWithBlogCount() ([]core.Tag, error)
	GetAllBlogs() ([]core.Blog, error)
	GetAllBlogsByTag(tagID string) ([]core.Blog, error)
}
