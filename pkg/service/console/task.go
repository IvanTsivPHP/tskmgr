package console

import (
	"tskmgr/pkg/storage"
)

type TaskService struct {
	TaskStorage storage.TaskStorage
}

func NewTaskService(storage storage.TaskStorage) *TaskService {
	return &TaskService{
		TaskStorage: storage,
	}
}
