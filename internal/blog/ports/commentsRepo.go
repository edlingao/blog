package ports

import "github.com/edlingao/internal/blog/core"

type CommentsRepository interface {
	SaveComment(comment *core.Comment) (*core.Comment, error)
	ReplyComment(parentCommentID string, childComment *core.Comment) (*core.Comment, error)
	GetCommentsByPostID(postID string) ([]*core.Comment, error)
	DeleteComment(commentID string) error
	GetCommentByID(commentID string) (*core.Comment, error)
	GetAllCommentsByPostIDWithChildren(postID string) ([]*core.Comment, error)
}
