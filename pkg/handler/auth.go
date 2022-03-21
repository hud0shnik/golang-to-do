package handler

import (
	"net/http"
	"todo-app"

	"github.com/gin-gonic/gin"
)

// Хендлер для эндпоинта регистрации
func (h *Handler) signUp(c *gin.Context) {

	// Структура, в которую будет производиться запись данных из JSON от пользователей
	var input todo.User

	// Распарсивает данные из JSON в структуру input
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Создание пользователя
	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Запись в ответ статус кода 200 и id пользователя
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})

}

// Структура входных данных для авторизации
type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Хендлер для эндпоинта авторизации
func (h *Handler) signIn(c *gin.Context) {

	// Структура ввода для авторизации
	var input signInInput

	// Запись значений из контекста
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Генерация токена
	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Запись в ответ статус кода 200 и токена
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
