package service

import (
	"errors"
	"fmt"
	"time"

	"forum/models"
	"forum/repository"
)

// GetAllPosts from posts and likes tables
func GetAllPosts(repos *repository.Repository, currentUserId int) ([]models.Post, error) {
	allPosts, err := repos.Posts.GetAllPosts(currentUserId)
	if err != nil {
		return allPosts, errors.New("can't get all posts")
	}
	return allPosts, nil
}

// GetPostById from posts, comments and likes tables
func GetPostById(repos *repository.Repository, id int) (models.Post, error) {
	post, err := repos.Posts.GetPostById(id)
	if err != nil {
		return post, errors.New("can't get all posts")
	}
	return post, nil
}

// GetPostById from posts, comments and likes tables
func DeletePostById(repos *repository.Repository, id int) error {
	err := repos.Posts.DeletePostById(id)
	if err != nil {
		return errors.New("can't get all posts")
	}
	return nil
}

// AddPost to posts table
func AddPost(repos *repository.Repository, post models.Post, categories []int, currentUserId int) (int, error) {
	post.Date = time.Now()
	id, err := repos.Posts.CreatePost(post, currentUserId)
	if err != nil {
		return 0, fmt.Errorf("DB can't add post: %w", err)
	}
	for _, catId := range categories {
		if err := repos.Posts.AddCategoryToPost(id, catId); err != nil {
			return 0, fmt.Errorf("DB can't add category: %w", err)
		}
	}
	return id, nil
}

// AddPost to posts table
func SetPost(repos *repository.Repository, post models.Post, categories []int, currentUserId int) (models.Post, error) {
	if repos.Posts.DeleteCategoriesToPost(post.Id) != nil {
		return post, fmt.Errorf("DB can't delete cats: %w", repos.Posts.DeleteCategoriesToPost(post.Id))
	}
	editedPost, err := repos.Posts.SetPost(post, currentUserId)
	if err != nil {
		return editedPost, fmt.Errorf("DB can't add post: %w", err)
	}

	for _, catId := range categories {
		if err := repos.Posts.AddCategoryToPost(post.Id, catId); err != nil {
			return editedPost, fmt.Errorf("DB can't add category: %w", err)
		}
	}
	return editedPost, nil
}
