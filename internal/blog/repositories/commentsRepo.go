package repositories

import (
	"strconv"

	"github.com/edlingao/internal/blog/core"
	"github.com/edlingao/internal/blog/queries"
	"github.com/jmoiron/sqlx"
)

type CommentsRepo struct {
	db *sqlx.DB
}

func NewCommentsRepo(db *sqlx.DB) *CommentsRepo {
	return &CommentsRepo{
		db: db,
	}
}

func (commentsRepo *CommentsRepo) SaveComment(comment *core.Comment) (*core.Comment, error) {
	row, err := commentsRepo.db.NamedExec(queries.InsertComment, comment)
	if err != nil {
		return &core.Comment{}, err
	}

	commentID, err := row.LastInsertId()
	if err != nil {
		return &core.Comment{}, err
	}

	savedComment, err := commentsRepo.GetCommentByID(strconv.Itoa(int(commentID)))
	if err != nil {
		return &core.Comment{}, err
	}

	return savedComment, nil
}

func (commentsRepo *CommentsRepo) ReplyComment(parentCommentID string, childComment *core.Comment) (*core.Comment, error) {
	row, err := commentsRepo.db.NamedExec(queries.ReplyComment, childComment)
	if err != nil {
		return &core.Comment{}, err
	}

	commentID, err := row.LastInsertId()
	if err != nil {
		return &core.Comment{}, err
	}

	savedComment, err := commentsRepo.GetCommentByID(strconv.Itoa(int(commentID)))
	if err != nil {
		return &core.Comment{}, err
	}

	return savedComment, nil
}

func (commentsRepo *CommentsRepo) GetCommentsByPostID(postID string) ([]*core.Comment, error) {
	var comments []*core.Comment
	err := commentsRepo.db.Select(&comments, queries.GetCommentsByPostID, postID)
	if err != nil {
		return []*core.Comment{}, err
	}

	return comments, nil
}

func (commentsRepo *CommentsRepo) DeleteComment(commentID string) error {
	_, err := commentsRepo.db.Exec(queries.DeleteComment, commentID)
	if err != nil {
		return err
	}

	return nil
}

func (commentsRepo *CommentsRepo) GetCommentByID(commentID string) (*core.Comment, error) {
	comment := core.NewComment("", "", "")
	err := commentsRepo.db.Get(comment, queries.GetCommentByID, commentID)
	if err != nil {
		return &core.Comment{}, err
	}

	return comment, nil
}

func (commentsRepo *CommentsRepo) GetAllCommentsByPostIDWithChildren(postID string) ([]*core.Comment, error) {
	var comments []*core.Comment
	err := commentsRepo.db.Select(&comments, queries.GetCommentsByPostID, postID)
	if err != nil {
		return []*core.Comment{}, err
	}

	commentTree := core.BuildCommentTree(comments)
	return commentTree, nil
}

