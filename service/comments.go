package service

import (
	"errors"
	"fmt"
	"time"

	"forum/models"
	"forum/repository"
)

// GetAllComments from comments and likes tables
func GetAllComments(repos *repository.Repository, postId int) ([]models.Comment, error) {
	comments, err := repos.Comments.GetCommentsByPostId(postId)
	if err != nil {
		return comments, errors.New("can't get all comments")
	}
	return comments, nil
}

// AddComment to comments table
func AddComment(repos *repository.Repository, comm models.Comment, currentUserId int) (int, error) {
	comm.Date = time.Now()
	id, err := repos.Comments.CreateComment(comm, currentUserId)
	if err != nil {
		return 0, fmt.Errorf("DB can't add token: %w", err)
	}
	return id, nil
}

func GetCommentById(repos *repository.Repository, id int) (models.Comment, error) {
	comment, err := repos.Comments.GetCommentById(id)
	if err != nil {
		return comment, fmt.Errorf("Can't get comment: %w", err)
	}
	return comment, nil
}
func DeleteCommentById(repos *repository.Repository, id int) error {
	err := repos.Comments.DeleteCommentById(id)
	if err != nil {
		return fmt.Errorf("Can't get comment: %w", err)
	}
	return nil
}

func EditComment(repos *repository.Repository, comment models.Comment, currentUserId int) (models.Comment, error) {
	comment, err := repos.Comments.SetComment(comment, currentUserId)
	if err != nil {
		return comment, fmt.Errorf("Can't set comment: %w", err)
	}
	return comment, nil
}
