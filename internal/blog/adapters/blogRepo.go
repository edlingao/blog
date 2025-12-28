package adapters

import (
	"errors"
	"slices"

	"github.com/edlingao/internal/blog/core"
	"github.com/edlingao/internal/blog/queries"
	"github.com/jmoiron/sqlx"
)

const (
	FailedToSaveBlog       = "failed to save blog post"
	FailedToUpdateBlog     = "failed to update blog post"
	FailedToGetBlogByTitle = "failed to get blog ( does not exist )"
)

type BlogRepo struct {
	dbConnection *sqlx.DB
}

func NewBlogRepo(db *sqlx.DB) *BlogRepo {
	return &BlogRepo{
		dbConnection: db,
	}
}

func (blogRepo *BlogRepo) Save(blog *core.Blog) (*core.Blog, error) {
	_, err := blogRepo.dbConnection.NamedExec(queries.InsertPost, blog)
	if err != nil {
		return &core.Blog{}, errors.New(FailedToSaveBlog + err.Error())
	}

	blog, error := blogRepo.GetByTitle(blog.Title)
	if error != nil {
		return &core.Blog{}, errors.New(FailedToSaveBlog + error.Error())
	}

	return blog, nil
}

func (BlogRepo *BlogRepo) Update(blog *core.Blog) (*core.Blog, error) {
	_, err := BlogRepo.dbConnection.NamedExec(queries.UpdatePost, blog)
	if err != nil {
		return &core.Blog{}, errors.New(FailedToUpdateBlog + err.Error())
	}

	updatedBlog, error := BlogRepo.GetByTitle(blog.Title)
	if error != nil {
		return &core.Blog{}, errors.New(FailedToUpdateBlog + error.Error())
	}

	return updatedBlog, nil
}

func (blogRepo *BlogRepo) GetByTitle(postTitle string) (*core.Blog, error) {
	blog := core.NewBlog(postTitle)
	err := blogRepo.dbConnection.Get(blog, queries.GetByTitle, postTitle)
	if err != nil {
		return &core.Blog{}, errors.New(FailedToGetBlogByTitle + err.Error())
	}
	return blog, nil
}

func (blogRepo *BlogRepo) AddTagsToBlog(blogID string, tags []string) error {
	existingTags, err := blogRepo.GetTagsByBlogID(blogID)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		tag := &core.Tag{
			Name: tag,
		}

		err := blogRepo.dbConnection.Get(tag, queries.GetTagIDByName, tag.Name)
		if err != nil {
			return err
		}

		if slices.Contains(existingTags, *tag) {
			existingTags = slices.DeleteFunc(existingTags, func(t core.Tag) bool {
				return t.ID == tag.ID
			})
			continue
		}

		tagBlog := &core.PostTag{
			PostID: blogID,
			TagID:  tag.ID,
		}
		_, err = blogRepo.dbConnection.NamedExec(queries.AddTagToBlog, tagBlog)
		if err != nil {
			return err
		}
	}

	for _, tag := range existingTags {
		_, err := blogRepo.dbConnection.NamedExec(
			queries.RemoveTagsFromBlog,
			&core.PostTag{
				PostID: blogID,
				TagID:  tag.ID,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (blogRepo *BlogRepo) RemoveTagFromBlog(blogID string, tagID string) error {
	_, err := blogRepo.dbConnection.Exec(queries.RemoveTagsFromBlog, blogID, tagID)
	if err != nil {
		return err
	}

	return nil
}

func (blogRepo *BlogRepo) GetTagsByBlogID(blogID string) ([]core.Tag, error) {
	var tags []core.Tag
	err := blogRepo.dbConnection.Select(&tags, queries.GetTagsByBlogID, blogID)
	if err != nil {
		return []core.Tag{}, err
	}

	return tags, nil
}
