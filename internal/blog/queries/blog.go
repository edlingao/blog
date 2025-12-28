package queries

const GetByTitle = `
	SELECT id, title, url, description, md_url, reactions FROM posts WHERE title = :title;
`

const InsertPost = `
	INSERT INTO posts (title, url, description, md_url, reactions)
	VALUES (:title, :url, :description, :md_url, :reactions);
`

const UpdatePost = `
	UPDATE posts
	SET title = :title, url = :url, description = :description, md_url = :md_url, reactions = :reactions
	WHERE id = :id;
`

const AddTagToBlog = `
	INSERT OR IGNORE INTO post_tags (post_id, tag_id)
	VALUES (:post_id, :tag_id);
`

const GetTagIDByName = `
	SELECT id, name FROM tags WHERE name = :name;
`

const RemoveTagsFromBlog = `
	DELETE FROM post_tags WHERE post_id = :post_id AND tag_id = :tag_id;
`

const GetTagsByBlogID = `
	SELECT t.id, t.name
	FROM tags t
	INNER JOIN post_tags pt ON t.id = pt.tag_id
	WHERE pt.post_id = :post_id;
`
