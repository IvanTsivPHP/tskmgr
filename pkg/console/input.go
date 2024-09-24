package console

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"tskmgr/pkg/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

func HandleDatabaseMenu() {
	formatter := &SimpleErrorFormatter{}
	for {
		ShowDatabaseMenu()

		dbChoice, err := ReadIntInput("Введите номер действия: ")
		if err != nil {
			HandleError(err, formatter)
		}

		switch dbChoice {
		case 1:
			fmt.Println("Запуск миграций вверх...")
			storage.RunMigrations("up") // Вызов функции для миграции вверх
		case 2:
			fmt.Println("Запуск миграций вниз...")
			storage.RunMigrations("down") // Вызов функции для миграции вниз
		case 3:
			// Вернуться в главное меню
			return
		default:
			fmt.Println("Некорректный выбор. Попробуйте снова.")
		}
	}
}

func HandleTaskMenu(db *pgxpool.Pool) {
	taskStorage := &storage.TaskStorage{DB: db}

	for {
		ShowTaskMenu()

		taskChoice, err := ReadIntInput("Введите номер действия: ")

		if err != nil {
			fmt.Println("Ошибка ввода. Пожалуйста, введите корректное число.")
			continue
		}

		switch taskChoice {
		case 1:
			fmt.Println("Создание новой задачи...")
			task := storage.Task{}
			// Пример ввода данных задачи
			fmt.Print("Введите название задачи: ")
			fmt.Scan(&task.Title)
			fmt.Print("Введите содержимое задачи: ")
			fmt.Scan(&task.Content)
			fmt.Print("Введите ID автора задачи: ")
			fmt.Scan(&task.AuthorID)
			fmt.Print("Введите ID ответственного за задачу: ")
			fmt.Scan(&task.AssignedID)

			fmt.Println("Введите метки через запятую (или оставьте пустым): ")
			var labelsInput string
			fmt.Scanln(&labelsInput)
			if labelsInput != "" {
				task.Labels = strings.Split(labelsInput, ",")
			}
			// Создание задачи
			taskID, err := taskStorage.CreateTask(task)
			if err != nil {
				fmt.Println("Ошибка создания задачи:", err)
			} else {
				fmt.Printf("Задача создана с ID %d\n", taskID)
			}
		case 2:
			fmt.Println("Получение списка всех задач...")
			tasks, err := taskStorage.GetAllTasks()
			if err != nil {
				fmt.Println("Ошибка получения задач:", err)
			} else {
				tasksTable(tasks)

			}

		case 3:
			fmt.Println("Получение задач по автору...")
			fmt.Print("Введите ID автора: ")
			var authorID int
			fmt.Scan(&authorID)
			tasks, err := taskStorage.GetTasksByAuthor(authorID)
			if err != nil {
				fmt.Println("Ошибка получения задач:", err)
			} else {
				for _, task := range tasks {
					fmt.Printf("ID: %d, Title: %s, Content: %s\n", task.ID, task.Title, task.Content)
				}
			}
		case 4:
			fmt.Println("Получение задач по метке...")
			fmt.Print("Введите метку: ")
			var label string
			fmt.Scan(&label)
			tasks, err := taskStorage.GetTasksByLabel(label)
			if err != nil {
				fmt.Println("Ошибка получения задач:", err)
			} else {
				for _, task := range tasks {
					fmt.Printf("ID: %d, Title: %s, Content: %s\n", task.ID, task.Title, task.Content)
				}
			}
		case 5:
			fmt.Println("Обновление задачи по ID...")
			fmt.Print("Введите ID задачи: ")
			task := storage.Task{}
			fmt.Scan(&task.ID)
			fmt.Print("Введите новое название задачи: ")
			fmt.Scanln(&task.Title)
			fmt.Print("Введите новое содержимое задачи: ")
			fmt.Scanln(&task.Content)
			fmt.Print("Введите ID нового автора задачи: ")
			fmt.Scanln(&task.AuthorID)
			fmt.Print("Введите ID нового ответственного за задачу: ")
			fmt.Scanln(&task.AssignedID)
			err := taskStorage.UpdateTask(task)
			if err != nil {
				fmt.Println("Ошибка обновления задачи:", err)
			} else {
				fmt.Println("Задача обновлена.")
			}
		case 6:
			fmt.Println("Удаление задачи по ID...")
			fmt.Print("Введите ID задачи: ")
			var taskID int
			fmt.Scan(&taskID)
			err := taskStorage.DeleteTask(taskID)
			if err != nil {
				fmt.Println("Ошибка удаления задачи:", err)
			} else {
				fmt.Println("Задача удалена.")
			}
		case 7:
			// Вернуться в главное меню
			return
		default:
			fmt.Println("Некорректный выбор. Попробуйте снова.")
		}
	}
}

func HandleMainMenu() {
	formatter := &DetailedErrorFormatter{}
	db, err := storage.NewDBConnection() // Инициализация подключения к базе данных
	if err != nil {
		fmt.Println("Ошибка подключения к базе данных:", err)
		return
	}
	defer db.Close() // Закрытие подключения при завершении работы

	for {
		ShowMainMenu()

		var mainChoice int

		mainChoice, err := ReadIntInput("Введите номер действия: ")
		if err != nil {
			HandleError(err, formatter)
			continue
		}

		switch mainChoice {
		case 1:
			HandleTaskMenu(db)
		case 2:
			HandleDatabaseMenu()
		case 3:
			fmt.Println("Выход из программы.")
			os.Exit(0)
		default:
			fmt.Println("Некорректный выбор. Попробуйте снова.")
		}
	}
}

// ReadIntInput считывает целое число с консоли и возвращает его вместе с ошибкой (если таковая имеется).
func ReadIntInput(prompt string) (int, error) {

	fmt.Print(prompt)
	var input string
	_, err := fmt.Scan(&input)

	if err != nil {
		return 0, err
	}
	// Преобразуем строку в целое число
	value, err := strconv.Atoi(input)

	return value, err
}

// ReadStringInput считывает строку с консоли и возвращает её.
func ReadStringInput(prompt string) (string, error) {
	fmt.Print(prompt)
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return "", err
	}
	return input, nil
}
