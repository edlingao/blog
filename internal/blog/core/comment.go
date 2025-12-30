package core

type Comment struct {
	ID        string     `db:"id" json:"id"`
	PostID    string     `db:"post_id" json:"post_id"`
	CommentID *string    `db:"comment_id" json:"comment_id"`
	Author    string     `db:"author" json:"author"`
	Reactions string     `db:"reactions" json:"reactions"`
	Content   string     `db:"content" json:"content"`
	Children  []*Comment `db:"children,omitempty" json:"children"`
	CreatedAt string     `db:"created_at" json:"created_at"`
	UpdatedAt string     `db:"updated_at" json:"updated_at"`
}

func NewComment(postID, author, content string) *Comment {
	return &Comment{
		PostID:  postID,
		Author:  author,
		Content: content,
	}
}

func (comment *Comment) AddParentComment(parentID string) {
	comment.CommentID = &parentID
}

func (comment *Comment) ReplyComment(child *Comment) {
	comment.Children = append(comment.Children, child)
}

func (comment *Comment) EditComment(newContent string) {
	comment.Content = newContent
}

func BuildCommentTree(comments []*Comment) []*Comment {
	byID := make(map[string]*Comment)
	var roots []*Comment

	for _, c := range comments {
		c.Children = []*Comment{}
		byID[c.ID] = c
	}

	for _, c := range comments {
		if c.CommentID == nil {
			roots = append(roots, c)
		} else if parent, ok := byID[*c.CommentID]; ok {
			parent.Children = append(parent.Children, c)
		}
	}
	return roots
}
