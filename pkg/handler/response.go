package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Структура ошибки
type errorResponse struct {
	Message string `json:"message"`
}

// Структура статуса приложения
type statusResponse struct {
	Status string `json:"status"`
}

// Функция для обработки ошибок
func newErrorResponse(c *gin.Context, statusCode int, message string) {

	// Вывод ошибки в консоль
	logrus.Error(message)

	// Метод, который блокирует выполнение последующих обработчиков
	// и записывает в ответ статускод и текст ошибки
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
