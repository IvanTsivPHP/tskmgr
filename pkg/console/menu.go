package console

import (
	"fmt"
)

func ShowMainMenu() {
	fmt.Println("Главное меню:")
	fmt.Println("1. Управление задачами")
	fmt.Println("2. База данных")
	fmt.Println("3. Выход")
}

func ShowDatabaseMenu() {
	fmt.Println("Меню базы данных:")
	fmt.Println("1. Выполнить миграции вверх (up)")
	fmt.Println("2. Выполнить миграции вниз (down)")
	fmt.Println("3. Вернуться в главное меню")
}

func ShowTaskMenu() {
	fmt.Println("Меню задач:")
	fmt.Println("1. Создать новую задачу")
	fmt.Println("2. Получить список всех задач")
	fmt.Println("3. Получить задачи по автору")
	fmt.Println("4. Получить задачи по метке")
	fmt.Println("5. Обновить задачу по ID")
	fmt.Println("6. Удалить задачу по ID")
	fmt.Println("7. Вернуться в главное меню")
}
