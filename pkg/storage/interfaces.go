package storage

import "tskmgr/pkg/models"

type TaskStorageInterface interface {
	CreateTask(task models.Task) (int, error)
	GetAllTasks() ([]models.Task, error)
	GetTasksByAuthor(authorID int) ([]models.Task, error)
	GetTasksByLabel(labelName string) ([]models.Task, error)
	UpdateTask(task models.Task) error
	DeleteTask(id int) error
}
