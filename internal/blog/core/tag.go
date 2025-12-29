package core

type Tag struct {
	ID    string `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Emoji string `json:"emoji" db:"emoji"`
	Count int    `json:"count" db:"count,omitempty"`
}

type PostTag struct {
	PostID string `json:"post_id" db:"post_id"`
	TagID  string `json:"tag_id" db:"tag_id"`
}
