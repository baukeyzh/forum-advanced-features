package service

import (
	"fmt"

	"forum/models"
	"forum/repository"
)

// GetAllComments from comments and likes tables
func AddLikePost(repos *repository.Repository, like models.LikePost, currentUserId int) (int, error) {
	id, err := repos.Likes.AddLikePost(like, currentUserId)
	if err != nil {
		return 0, fmt.Errorf("DB can't add like to post: %w", err)
	}
	return id, nil
}

// AddComment to comments table
func AddLikeComment(repos *repository.Repository, like models.LikeComment, currentUserId int) (int, error) {
	id, err := repos.Likes.AddLikeComment(like, currentUserId)
	if err != nil {
		return 0, fmt.Errorf("DB can't add like to comment: %w", err)
	}
	return id, nil
}
