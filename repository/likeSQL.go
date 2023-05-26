package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"forum/models"

	_ "github.com/mattn/go-sqlite3"
)

type likeSQL struct {
	db *sql.DB
}

// NewlikeSQL create new database struct
func NewlikeSQL(db *sql.DB) *likeSQL {
	return &likeSQL{db: db}
}

// AddLikePost
// INSERT INTO posts_likes (user_id, post_id, type) values (1,2,true)
func (r *likeSQL) AddLikePost(like models.LikePost, currentUserId int) (int, error) {
	var id int
	postId, err := strconv.Atoi(like.PostID)
	if err != nil {
		return 0, err
	}
	likeFromDb, err := r.GetLikeByPostUser(postId, like.AuthorID)
	postAuthorId, err := GetAuthorId(r, 0, postId)
	query := ""
	typeOfAction := ""
	if like.Type == true {
		typeOfAction = "like"
	} else {
		typeOfAction = "dislike"
	}
	if likeFromDb.Id != 0 {
		if likeFromDb.Type != like.Type {
			query = fmt.Sprintf(`UPDATE %s SET type = ?  WHERE id = ?`, postsLikesTable)
			if _, err := r.db.Exec(query, like.Type, likeFromDb.Id); err != nil {
				return 0, fmt.Errorf("can't set like type: %w", err)
			}
			err = InsertLikeActivity(r, currentUserId, postId, 0, postAuthorId, typeOfAction)
			if err != nil {
				return 0, err
			}
			return likeFromDb.Id, nil
		} else {
			query = fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, postsLikesTable)
			if _, err := r.db.Exec(query, likeFromDb.Id); err != nil {
				return 0, fmt.Errorf("can't delete like: %w", err)
			}
			return 0, nil
		}
	} else {
		query = fmt.Sprintf(`INSERT INTO %s (user_id, post_id, type) values (?,?,?) RETURNING id`, postsLikesTable)
		row := r.db.QueryRow(query, like.AuthorID, like.PostID, like.Type)
		if err := row.Scan(&id); err != nil {
			return 0, err
		}
		err = InsertLikeActivity(r, currentUserId, postId, 0, postAuthorId, typeOfAction)
		if err != nil {
			return 0, err
		}
	}
	return id, nil
}

// AddLikeComment
// INSERT INTO comments_likes (user_id, comment_id, type) values (1,2,true)
func (r *likeSQL) AddLikeComment(like models.LikeComment, currentUserId int) (int, error) {
	var id int
	CommentID, err := strconv.Atoi(like.CommentID)
	if err != nil {
		return 0, err
	}
	likeFromDb, err := r.GetLikeByCommentUser(CommentID, like.AuthorID)
	commentAuthorId, err := GetAuthorId(r, CommentID, 0)
	if err != nil {
		return 0, err
	}
	query := ""
	typeOfAction := ""
	if like.Type == true {
		typeOfAction = "like"
	} else {
		typeOfAction = "dislike"
	}
	if likeFromDb.Id != 0 {
		if likeFromDb.Type != like.Type {
			query = fmt.Sprintf(`UPDATE %s SET type = ?  WHERE id = ?`, commentsLikesTable)
			if _, err := r.db.Exec(query, like.Type, likeFromDb.Id); err != nil {
				return 0, fmt.Errorf("can't set like type: %w", err)
			}
			err = InsertLikeActivity(r, currentUserId, 0, CommentID, commentAuthorId, typeOfAction)
			if err != nil {
				return 0, err
			}
			return likeFromDb.Id, nil
		} else {
			query = fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, commentsLikesTable)
			if _, err := r.db.Exec(query, likeFromDb.Id); err != nil {
				return 0, fmt.Errorf("can't delete like: %w", err)
			}
			return 0, nil
		}
	} else {
		query = fmt.Sprintf(`INSERT INTO %s (user_id, comment_id, type) values (?,?,?) RETURNING id`, commentsLikesTable)
		row := r.db.QueryRow(query, like.AuthorID, like.CommentID, like.Type)
		if err := row.Scan(&id); err != nil {
			return 0, err
		}
		err = InsertLikeActivity(r, currentUserId, 0, CommentID, commentAuthorId, typeOfAction)
		if err != nil {
			return 0, err
		}
	}
	return id, nil
}

// GetLikeByPostUser
// SELECT * FROM posts_likes WHERE post_id=1 AND user_id=2
func (r *likeSQL) GetLikeByPostUser(postId, userId int) (models.LikePost, error) {
	var like models.LikePost

	query := fmt.Sprintf("SELECT * FROM %s WHERE post_id= ? AND user_id= ?", postsLikesTable)
	err := r.db.QueryRow(query, postId, userId).Scan(
		&like.Id,
		&like.AuthorID,
		&like.PostID,
		&like.Type,
	)
	if err == sql.ErrNoRows {
		return like, err
	}
	if err != nil {
		return like, fmt.Errorf("can't get the like of this post: %w", err)
	}
	return like, nil
}

// GetLikeByCommentUser
// SELECT * FROM comments_likes WHERE post_id=1 AND  user_id=2
func (r *likeSQL) GetLikeByCommentUser(commentId, userId int) (models.LikeComment, error) {
	var like models.LikeComment

	query := fmt.Sprintf("SELECT * FROM %s WHERE comment_id= ? AND user_id=?", commentsLikesTable)
	err := r.db.QueryRow(query, commentId, userId).Scan(
		&like.Id,
		&like.AuthorID,
		&like.CommentID,
		&like.Type,
	)
	if err == sql.ErrNoRows {
		return like, nil
	}
	if err != nil {
		return like, fmt.Errorf("can't get all the like of this comment: %w", err)
	}
	return like, nil
}

func InsertLikeActivity(r *likeSQL, userId, postId, commentId, authorId int, action string) error {
	now := time.Now()
	query := fmt.Sprintf(`INSERT INTO %s (user_id, post_id, comment_id, author_id, action, is_shown, created_at) values (?,?,?,?,?,?,?) RETURNING id`, activityTable)
	if _, err := r.db.Exec(query, userId, postId, commentId, authorId, action, false, now); err != nil {
		return fmt.Errorf("can't set activity: %w", err)
	}
	return nil
}

func GetAuthorId(r *likeSQL, commentId, postId int) (int, error) {
	var authorId models.AuthorID
	var elemId int
	var table string
	if postId != 0 {
		table = postsTable
		elemId = postId
	} else {
		table = commentsTable
		elemId = commentId
	}
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
