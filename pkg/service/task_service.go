package service

import (
	"errors"
	"time"
	"tskmgr/pkg/models"
	"tskmgr/pkg/storage"
)

type TaskService struct {
	storage storage.TaskStorageInterface
}

func NewTaskService(repo storage.TaskStorageInterface) *TaskService {
	return &TaskService{storage: repo}
}

// Создание задачи с добавлением текущей метки времени
func (s *TaskService) CreateTask(task models.Task) (int, error) {
	// Добавление бизнес-логики: метка времени открытия задачи
	if task.Title == "" || task.Content == "" {
		return 0, errors.New("title and content cannot be empty")
	}
	task.Opened = time.Now().Unix()
	return s.storage.CreateTask(task)
}

// Получение всех задач
func (s *TaskService) GetAllTasks() ([]models.Task, error) {
	return s.storage.GetAllTasks()
}

// Получение задач по автору
func (s *TaskService) GetTasksByAuthor(authorID int) ([]models.Task, error) {
	if authorID == 0 {
		return nil, errors.New("author ID cannot be zero")
	}
	return s.storage.GetTasksByAuthor(authorID)
}

// Получение задач по метке
func (s *TaskService) GetTasksByLabel(label string) ([]models.Task, error) {
	if label == "" {
		return nil, errors.New("label cannot be empty")
	}
	return s.storage.GetTasksByLabel(label)
}

// Обновление задачи
func (s *TaskService) UpdateTask(task models.Task) error {
	if task.ID == 0 {
		return errors.New("task ID is required")
	}
	return s.storage.UpdateTask(task)
}

// Удаление задачи
func (s *TaskService) DeleteTask(id int) error {
	if id == 0 {
		return errors.New("task ID is required")
	}
	return s.storage.DeleteTask(id)
}
