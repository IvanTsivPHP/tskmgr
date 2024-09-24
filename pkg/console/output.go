package console

import (
	"fmt"
	"strings"
	"time"
	"tskmgr/pkg/storage"
)

func tasksTable(tasks []storage.Task) {
	//Формируем шапку
	fmt.Printf("%-5s | %-20s | %-20s | %-10s | %-10s | %-20s | %-40s | %-20s\n", "ID", "Открыта", "Закрыта", "Автор", "Исполнитель", "Заголовок", "Содержание", "Метки")
	fmt.Println(strings.Repeat("-", 150)) // Разделитель между шапкой и содержимым

	// Вывод каждой задачи
	for _, task := range tasks {

		openedTime := time.Unix(task.Opened, 0).Format("2006-01-02 15:04:05")
		closedTime := ""
		if task.Closed != 0 {
			closedTime = time.Unix(task.Closed, 0).Format("2006-01-02 15:04:05")
		}

		labelsString := strings.Join(task.Labels, ", ") // Соединение меток в строку через запятую
		fmt.Printf(
			"%-5d | %-20s | %-20s | %-10d | %-11d | %-20s | %-40s | %-20s\n",
			task.ID, openedTime, closedTime, task.AuthorID, task.AssignedID, task.Title, task.Content, labelsString)
	}
}
