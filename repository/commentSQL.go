package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"forum/models"

	_ "github.com/mattn/go-sqlite3"
)

type commentSQL struct {
	db *sql.DB
}

// New create new database.
func NewCommentSQL(db *sql.DB) *commentSQL {
	return &commentSQL{db: db}
}

// INSERT INTO comments (user_id, date, post_id , content) values (1, "2023-05-01 13:35:04.556898354+06:00" , 1, "golang top, js for girls");
// CreateComment
func (r *commentSQL) CreateComment(comment models.Comment, currentUserId int) (int, error) {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s (user_id, date, post_id, content) values (?,?,?,?) RETURNING id`, commentsTable)

	row := r.db.QueryRow(query, comment.AuthorID, comment.Date, comment.PostID, comment.Content)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	postAuthorId, err := GetPostAuthorId(r, comment.PostID)
	err = InsertCommentActivity(r, currentUserId, 0, id, postAuthorId, "create")
	if err != nil {
		return 0, err
	}
	return id, nil
}

// CreateComment
func (r *commentSQL) SetComment(comment models.Comment, currentUserId int) (models.Comment, error) {
	query := fmt.Sprintf(`UPDATE %s SET content = ? WHERE id = ?`, commentsTable)
	if _, err := r.db.Exec(query, comment.Content, comment.Id); err != nil {
		return comment, fmt.Errorf("can't set comment: %w", err)
	}
	err := InsertCommentActivity(r, currentUserId, 0, comment.Id, comment.AuthorID, "edit")
	if err != nil {
		return comment, err
	}
	return comment, nil
}
func (r *commentSQL) DeleteCommentById(commentId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, commentsTable)
	if _, err := r.db.Exec(query, commentId); err != nil {
		return fmt.Errorf("can't set comment: %w", err)
	}
	return nil
}

func (r *commentSQL) GetCommentsByPostId(postId int) ([]models.Comment, error) {
	var allComments []models.Comment

	query := fmt.Sprintf(`
	SELECT c.* , u.username as user_name,
		(SELECT Count (*) FROM %s cl WHERE cl.comment_id = c.id and type = true) as likes,
		(SELECT Count (*) FROM %s cl WHERE cl.comment_id = c.id and type = false) as dislikes
	FROM %s c 
	LEFT JOIN %s u ON u.id = c.user_id
	WHERE post_id=?;`, commentsLikesTable, commentsLikesTable, commentsTable, usersTable)
	rows, err := r.db.Query(query, postId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("no posts found: %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("can't get posts: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var comment models.Comment
		if err = rows.Scan(
			&comment.Id,
			&comment.AuthorID,
			&comment.Date,
			&comment.PostID,
			&comment.Content,
			&comment.AuthorName,
			&comment.Likes,
			&comment.Dislikes,
		); err != nil {
			return nil, fmt.Errorf("can't scan all comments: %w", err)
		}
		allComments = append(allComments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get all comments: %w", err)
	}
	return allComments, nil
}

func (r *commentSQL) GetCommentById(commentId int) (models.Comment, error) {
	var comment models.Comment
	query := fmt.Sprintf(`
		SELECT c.* , u.username as user_name,
			(SELECT Count (*) FROM %s cl WHERE cl.comment_id = c.id and type = true) as likes,
			(SELECT Count (*) FROM %s cl WHERE cl.comment_id = c.id and type = false) as dislikes
		FROM %s c 
		LEFT JOIN %s u ON u.id = c.user_id
		WHERE c.id=?;`, commentsLikesTable, commentsLikesTable, commentsTable, usersTable)
	err := r.db.QueryRow(query, commentId).Scan(
		&comment.Id,
		&comment.AuthorID,
		&comment.Date,
		&comment.PostID,
		&comment.Content,
		&comment.AuthorName,
		&comment.Likes,
		&comment.Dislikes,
	)
	if err != nil {
		return comment, fmt.Errorf("can't scan all comments: %w", err)
	}
	return comment, nil
}

func InsertCommentActivity(r *commentSQL, userId, postId, commentId, authorId int, action string) error {
	query := fmt.Sprintf(`INSERT INTO %s (user_id, post_id, comment_id, author_id, action, is_shown, created_at) values (?,?,?,?,?,?,?) RETURNING id`, activityTable)
	if _, err := r.db.Exec(query, userId, postId, commentId, authorId, action, false, time.Now()); err != nil {
		return fmt.Errorf("can't set activity: %w", err)
	}
	return nil
}

func GetPostAuthorId(r *commentSQL, postId string) (int, error) {
	var authorId models.AuthorID
	var elemId string
	var table string
	table = postsTable
	elemId = postId
	query := fmt.Sprintf(
		`SELECT user_id 
			FROM %s WHERE id= ?`, table)
	err := r.db.QueryRow(query, elemId).Scan(
		&authorId.AuthorId,
	)
	if err == sql.ErrNoRows {
		return authorId.AuthorId, nil
	}
	if err != nil {
		return authorId.AuthorId, fmt.Errorf("can't get all the like of this comment: %w", err)
	}
	return authorId.AuthorId, nil
}
