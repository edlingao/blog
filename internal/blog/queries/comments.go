package queries

const InsertComment = `
INSERT INTO comments (post_id, author, content, created_at)
VALUES (:post_id, :author, :content, date('now'));
`

const ReplyComment = `
INSERT INTO comments (post_id, author, content, comment_id, updated_at)
VALUES (:post_id, :author, :content, :comment_id, date('now'));
`

const GetCommentsByPostID = `
SELECT id, post_id, author, content, comment_id, created_at, updated_at
FROM comments
WHERE post_id = ?;
`

const DeleteComment = `
UPDATE comments
SET content = '[deleted]', updated_at = date('now')
WHERE id = ?;
`

const GetCommentByID = `
SELECT id, post_id, author, content, comment_id, created_at, updated_at
FROM comments
WHERE id = ?;
`

// Feed the children to a named 'children' field in the result
const GetCommentsByPostIDWithAllChildren = `
WITH RECURSIVE comment_tree AS (
		SELECT id, post_id, author, content, comment_id, created_at, updated_at
		FROM comments
		WHERE post_id = ? AND comment_id IS NULL
		UNION ALL
		SELECT c.id, c.post_id, c.author, c.content, c.comment_id, c.created_at, c.updated_at
		FROM comments c
		INNER JOIN comment_tree ct ON c.comment_id = ct.id
	)
	SELECT *
	FROM comment_tree;
`
