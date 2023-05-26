package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/models"

	_ "github.com/mattn/go-sqlite3"
)

type activitySQL struct {
	db *sql.DB
}

// New create new database.
func NewActivitySQL(db *sql.DB) *activitySQL {
	return &activitySQL{db: db}
}

// SELECT a.* , cu.username as user_name, au.username as author_name FROM activity a LEFT JOIN users cu ON cu.id = a.user_id LEFT JOIN users au ON au.id = a.author_id WHERE user_id=1 OR author_id=1
// GetActivitysByPostId
func (r *activitySQL) GetActivitiesByCurrentUserId(currentUserId int) ([]models.Activity, error) {
	var allActivities []models.Activity
	query := fmt.Sprintf(`
	SELECT a.* , cu.username as user_name, au.username as author_name,
	COALESCE((SELECT c.post_id from %s c where c.id=a.comment_id), 0) comment_post_id
	FROM %s a
	LEFT JOIN %s cu ON cu.id = a.user_id
	LEFT JOIN %s au ON au.id = a.author_id
	WHERE author_id=?`, commentsTable, activityTable, usersTable, usersTable)
	rows, err := r.db.Query(query, currentUserId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("no activities found: %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("can't get activities: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var activity models.Activity
		if err = rows.Scan(
			&activity.Id,
			&activity.UserId,
			&activity.CommentId,
			&activity.PostId,
			&activity.AuthorId,
			&activity.IsShown,
			&activity.Action,
			&activity.CreatedAt,
			&activity.UserName,
			&activity.AuthorName,
			&activity.CommentPostId,
		); err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("can't scan all activitys: %w", err)
		}
		allActivities = append(allActivities, activity)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get all activitys: %w", err)
	}
	err = SetActivitiesShown(r, currentUserId)
	if err != nil {
		return allActivities, err
	}
	return allActivities, nil
}

func (r *activitySQL) GetActivitiesCountByCurrentUserId(currentUserId int) (int, error) {
	var activitiesCount int
	query := fmt.Sprintf(`
	SELECT count(*)
	FROM %s a
	WHERE author_id=? AND is_shown=false`, activityTable)
	row := r.db.QueryRow(query, currentUserId)
	if err := row.Scan(&activitiesCount); err != nil {
		return 0, err
	}
	return activitiesCount, nil
}
func SetActivitiesShown(r *activitySQL, currentUserId int) error {
	query := fmt.Sprintf(`UPDATE %s SET is_shown=? WHERE author_id=?`, activityTable)
	if _, err := r.db.Exec(query, true, currentUserId); err != nil {
		return fmt.Errorf("can't set activity: %w", err)
	}
	return nil
}
