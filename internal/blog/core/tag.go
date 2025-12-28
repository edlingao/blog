package core

type Tag struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type PostTag struct {
	PostID string `json:"post_id" db:"post_id"`
	TagID  string `json:"tag_id" db:"tag_id"`
}


