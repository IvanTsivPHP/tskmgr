package dto

import (
	"fmt"
	"strings"
	"time"
	"tskmgr/pkg/models"
)

type TaskConsoleDTO struct {
	Title      string
	Content    string
	AuthorID   int
	AssignedID int
	Labels     string // Строка меток, разделенных запятыми
	Opened     string // Открытое время в формате YYYY-MM-DD HH:MM:SS
	Closed     string // Закрытое время в формате YYYY-MM-DD HH:MM:SS (может быть пустым)
}

// Конвертация задачи из консольного формата в внутреннюю модель
func (t TaskConsoleDTO) ToModel() (models.Task, error) {
	labels := strings.Split(t.Labels, ",")
	for i := range labels {
		labels[i] = strings.TrimSpace(labels[i])
	}

	// Парсинг времени открытия и закрытия
	openedTime, err := time.Parse("2006-01-02 15:04:05", t.Opened)
	if err != nil {
		return models.Task{}, fmt.Errorf("ошибка парсинга времени открытия: %w", err)
	}

	var closedTime int64
	if t.Closed != "" {
		parsedClosed, err := time.Parse("2006-01-02 15:04:05", t.Closed)
		if err != nil {
			return models.Task{}, fmt.Errorf("ошибка парсинга времени закрытия: %w", err)
		}
		closedTime = parsedClosed.Unix()
	}

	return models.Task{
		Title:      t.Title,
		Content:    t.Content,
		AuthorID:   t.AuthorID,
		AssignedID: t.AssignedID,
		Labels:     labels,
		Opened:     openedTime.Unix(),
		Closed:     closedTime,
	}, nil
}

// Конвертация задачи из внутренней модели в консольный формат
func FromModelToConsole(task models.Task) TaskConsoleDTO {
	openedTime := time.Unix(task.Opened, 0).Format("2006-01-02 15:04:05")
	var closedTime string
	if task.Closed != 0 {
		closedTime = time.Unix(task.Closed, 0).Format("2006-01-02 15:04:05")
	}

	return TaskConsoleDTO{
		Title:      task.Title,
		Content:    task.Content,
		AuthorID:   task.AuthorID,
		AssignedID: task.AssignedID,
		Labels:     strings.Join(task.Labels, ", "),
		Opened:     openedTime,
		Closed:     closedTime,
	}
}
