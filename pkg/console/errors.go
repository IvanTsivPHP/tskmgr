package console

import (
	"fmt"
)

// ErrorFormatter форматирует сообщения об ошибках
type ErrorFormatter interface {
	Format(err error) string
}

// SimpleErrorFormatter просто возвращает строковое представление ошибки
type SimpleErrorFormatter struct{}

// Format возвращает строковое представление ошибки
func (f *SimpleErrorFormatter) Format(err error) string {
	return fmt.Sprintf("Ошибка: %s", err)
}

// DetailedErrorFormatter возвращает более подробное сообщение об ошибке
type DetailedErrorFormatter struct {
	Context string
}

// Format возвращает строковое представление ошибки с контекстом
func (f *DetailedErrorFormatter) Format(err error) string {
	return fmt.Sprintf("Ошибка в %s: %s", f.Context, err)
}

// HandleError обрабатывает ошибку и выводит её в консоль
func HandleError(err error, formatter ErrorFormatter) {
	if err != nil {
		fmt.Println(formatter.Format(err))
	}
}
