package main

import (
	"fmt"
	"log"
	"tskmgr/pkg/models"
	"tskmgr/pkg/storage"
)

func main() {
	// Устанавливаем соединение с базой данных
	dbPool, err := storage.NewDBConnection()
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer dbPool.Close()

	// Инициализируем хранилище задач
	taskStorage := storage.NewTaskStorage(dbPool)

	// 1. Создание новой задачи
	task := models.Task{
		AuthorID:   1,
		AssignedID: 2,
		Title:      "Создать демо приложение",
		Content:    "Создать простое демонстрационное приложение для работы с задачами",
		Labels:     []string{"dev", "demo"},
	}

	taskID, err := taskStorage.CreateTask(task)
	if err != nil {
		log.Fatalf("Ошибка создания задачи: %v", err)
	}
	fmt.Printf("Задача создана с ID: %d\n", taskID)

	// 2. Получение всех задач
	allTasks, err := taskStorage.GetAllTasks()
	if err != nil {
		log.Fatalf("Ошибка получения всех задач: %v", err)
	}
	fmt.Println("Все задачи:")
	for _, t := range allTasks {
		fmt.Printf("ID: %d, Title: %s, Labels: %v\n", t.ID, t.Title, t.Labels)
	}

	// 3. Получение задач по автору
	authorTasks, err := taskStorage.GetTasksByAuthor(1)
	if err != nil {
		log.Fatalf("Ошибка получения задач по автору: %v", err)
	}
	fmt.Println("Задачи автора с ID 1:")
	for _, t := range authorTasks {
		fmt.Printf("ID: %d, Title: %s, Labels: %v\n", t.ID, t.Title, t.Labels)
	}

	// 4. Получение задач по метке
	labelTasks, err := taskStorage.GetTasksByLabel("demo")
	if err != nil {
		log.Fatalf("Ошибка получения задач по метке: %v", err)
	}
	fmt.Println("Задачи с меткой 'demo':")
	for _, t := range labelTasks {
		fmt.Printf("ID: %d, Title: %s, Labels: %v\n", t.ID, t.Title, t.Labels)
	}

	// 5. Обновление задачи
	taskToUpdate := models.Task{
		ID:         taskID,
		Title:      "Обновленное название",
		Content:    "Обновленное описание задачи",
		AssignedID: 3,
	}
	err = taskStorage.UpdateTask(taskToUpdate)
	if err != nil {
		log.Fatalf("Ошибка обновления задачи: %v", err)
	}
	fmt.Printf("Задача с ID %d обновлена.\n", taskID)

	// 6. Удаление задачи
	err = taskStorage.DeleteTask(taskID)
	if err != nil {
		log.Fatalf("Ошибка удаления задачи: %v", err)
	}
	fmt.Printf("Задача с ID %d удалена.\n", taskID)
}
