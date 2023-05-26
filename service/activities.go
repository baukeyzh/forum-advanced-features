package service

import (
	"errors"
	"fmt"

	"forum/models"
	"forum/repository"
)

func GetActivitiesByCurrentUserId(repos *repository.Repository, currentUserId int) ([]models.Activity, error) {
	activities, err := repos.Activities.GetActivitiesByCurrentUserId(currentUserId)
	if err != nil {
		fmt.Println(err)
		return activities, errors.New("can't get all activities")
	}
	return activities, nil
}

func GetActivitiesCountByCurrentUserId(repos *repository.Repository, currentUserId int) (int, error) {
	activitiesCount, err := repos.Activities.GetActivitiesCountByCurrentUserId(currentUserId)
	if err != nil {
		fmt.Println(err)
		return activitiesCount, errors.New("can't get all activities")
	}
	return activitiesCount, nil
}
