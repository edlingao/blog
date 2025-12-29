package core

type Comment struct {
	ID        string    `db:"id" json:"id"`
	PostID    string    `db:"post_id" json:"post_id"`
	CommentID string    `db:"comment_id" json:"comment_id"`
	Reactions string    `db:"reactions" json:"reactions"`
	Content   string    `db:"content" json:"content"`
	Children  []Comment `db:"-" json:"children"`
}

func NewComment(postID, commentID, content string) *Comment {
	return &Comment{
		PostID:    postID,
		CommentID: commentID,
		Reactions: "",
		Content:   "",
	}
}
