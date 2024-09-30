package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"tskmgr/pkg/models"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Структура данных для задачи

// Хранилище для работы с задачами
type TaskStorage struct {
	DB *pgxpool.Pool
}

func NewTaskStorage(db *pgxpool.Pool) *TaskStorage {
	return &TaskStorage{DB: db}
}

func NewDBConnection() (*pgxpool.Pool, error) {
	// Чтение конфигурации из файла
	config, err := LoadConfig("config.json")
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	// Формирование строки подключения
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%v/%s?sslmode=%s",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName, config.DBSSLMode)

	// Создание пула соединений
	dbpool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Проверка соединения
	err = dbpool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки соединения: %w", err)
	}

	log.Println("Успешное подключение к базе данных")

	return dbpool, nil
}

// Создание новой задачи с метками
func (s *TaskStorage) CreateTask(task models.Task) (int, error) {
	// Начало транзакции
	tx, err := s.DB.Begin(context.Background())
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(context.Background())

	// Создание задачи
	taskID, err := s.createTaskWithTx(tx, task)
	if err != nil {
		return 0, err
	}

	// Добавление меток (если есть)
	if len(task.Labels) > 0 {
		err = s.addLabelsToTaskWithTx(tx, taskID, task.Labels)
		if err != nil {
			return 0, err
		}
	}

	// Коммит транзакции
	err = tx.Commit(context.Background())
	if err != nil {
		return 0, err
	}

	return taskID, nil
}

// Получение списка всех задач
func (s *TaskStorage) GetAllTasks() ([]models.Task, error) {
	query := `
		SELECT t.id, t.opened, t.closed, t.author_id, t.assigned_id, t.title, t.content, COALESCE(array_agg(label.name), '{}') as labels
		FROM tasks t
		LEFT JOIN tasks_labels tl ON t.id = tl.task_id
		LEFT JOIN labels label ON tl.label_id = label.id
		GROUP BY t.id;
	`

	rows, err := s.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return s.scanTasks(rows)
}

// Получение задач по автору
func (s *TaskStorage) GetTasksByAuthor(authorID int) ([]models.Task, error) {
	query := `
		SELECT t.id, t.opened, t.closed, t.author_id, t.assigned_id, t.title, t.content, COALESCE(array_agg(label.name), '{}') as labels
		FROM tasks t
		LEFT JOIN tasks_labels tl ON t.id = tl.task_id
		LEFT JOIN labels label ON tl.label_id = label.id
		WHERE t.author_id = $1
		GROUP BY t.id;
	`

	rows, err := s.DB.Query(context.Background(), query, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.scanTasks(rows)
}

// Получение задач по метке
func (s *TaskStorage) GetTasksByLabel(labelName string) ([]models.Task, error) {
	query := `
		SELECT t.id, t.opened, t.closed, t.author_id, t.assigned_id, t.title, t.content, COALESCE(array_agg(label.name), '{}') as labels
		FROM tasks t
		JOIN tasks_labels tl ON t.id = tl.task_id
		JOIN labels label ON tl.label_id = label.id
		WHERE label.name = $1
		GROUP BY t.id;
	`

	rows, err := s.DB.Query(context.Background(), query, labelName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.scanTasks(rows)
}

// Обновление задачи по ID
func (s *TaskStorage) UpdateTask(task models.Task) error {
	query := "UPDATE tasks SET "
	var updates []string
	var args []interface{}
	argID := 1

	// Добавляем только те поля, которые переданы
	if task.Title != "" {
		updates = append(updates, fmt.Sprintf("title = $%d", argID))
		args = append(args, task.Title)
		argID++
	}
	if task.Content != "" {
		updates = append(updates, fmt.Sprintf("content = $%d", argID))
		args = append(args, task.Content)
		argID++
	}
	if task.AuthorID != 0 {
		updates = append(updates, fmt.Sprintf("author_id = $%d", argID))
		args = append(args, task.AuthorID)
		argID++
	}
	if task.AssignedID != 0 {
		updates = append(updates, fmt.Sprintf("assigned_id = $%d", argID))
		args = append(args, task.AssignedID)
		argID++
	}
	if task.Closed != 0 { // Проверка на наличие флага
		updates = append(updates, fmt.Sprintf("closed = $%d", argID))
		args = append(args, task.Closed)
		argID++
	}

	// Если нет изменений, возвращаем ошибку
	if len(updates) == 0 {
		return errors.New("no fields to update")
	}

	// Добавляем условие WHERE
	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d", argID)
	args = append(args, task.ID)

	// Выполняем обновление
	cmdTag, err := s.DB.Exec(context.Background(), query, args...)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("no task found with the specified ID")
	}

	return nil
}

// Удаление задачи по ID
func (s *TaskStorage) DeleteTask(id int) error {
	tx, err := s.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// Удаление связанных записей из tasks_labels
	_, err = tx.Exec(context.Background(), "DELETE FROM tasks_labels WHERE task_id = $1;", id)
	if err != nil {
		return err
	}

	// Удаление задачи
	cmdTag, err := tx.Exec(context.Background(), "DELETE FROM tasks WHERE id = $1;", id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("no task found with the specified ID")
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

// Вспомогательная функция для добавления меток к задаче
func (s *TaskStorage) addLabelsToTask(tx pgx.Tx, taskID int, labels []string) error {
	for _, label := range labels {
		var labelID int
		err := tx.QueryRow(context.Background(), `INSERT INTO labels (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name RETURNING id;`, label).Scan(&labelID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(context.Background(), `INSERT INTO tasks_labels (task_id, label_id) VALUES ($1, $2);`, taskID, labelID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *TaskStorage) addLabelsToTaskWithTx(tx pgx.Tx, taskID int, labels []string) error {
	for _, label := range labels {
		var labelID int
		// Вставка метки или обновление существующей
		err := tx.QueryRow(context.Background(), `
            INSERT INTO labels (name) VALUES ($1)
            ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name
            RETURNING id;
        `, label).Scan(&labelID)
		if err != nil {
			return err
		}

		// Привязка метки к задаче
		_, err = tx.Exec(context.Background(), `
            INSERT INTO tasks_labels (task_id, label_id) VALUES ($1, $2);
        `, taskID, labelID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *TaskStorage) createTaskWithTx(tx pgx.Tx, task models.Task) (int, error) {
	var taskID int
	task.Opened = time.Now().Unix()

	// Вставка задачи
	query := `
        INSERT INTO tasks (opened, closed, author_id, assigned_id, title, content)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id;
    `
	err := tx.QueryRow(context.Background(), query, task.Opened, task.Closed, task.AuthorID, task.AssignedID, task.Title, task.Content).Scan(&taskID)
	if err != nil {
		return 0, err
	}

	return taskID, nil
}

// Вспомогательная функция для обработки результата выборки задач
func (s *TaskStorage) scanTasks(rows pgx.Rows) ([]models.Task, error) {
	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		var labels sql.NullString
		err := rows.Scan(&task.ID, &task.Opened, &task.Closed, &task.AuthorID, &task.AssignedID, &task.Title, &task.Content, &labels)
		if err != nil {
			return nil, err
		}
		if labels.Valid && labels.String != "{}" {
			cleanedLabels := strings.Trim(labels.String, "{}")
			labelSlice := strings.Split(cleanedLabels, ", ")
			task.Labels = labelSlice
		} else {
			task.Labels = []string{}
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
