package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

// Метод идентификации пользователя
func (h *Handler) userIdentity(c *gin.Context) {

	// Получение значения из хедера авторизации
	header := c.GetHeader(authorizationHeader)

	// Валидация хедера авторизации
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 /*|| headerParts[0] != "Bearer" */ {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	// Парсинг токена
	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Запись id пользователя в контекст
	c.Set(userCtx, userId)
}

// Функция получения id пользователя из контекста
func getUserId(c *gin.Context) (int, error) {

	// Получение пользователя из контекста
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id is not found")
		return 0, errors.New("user id not found")
	}

	// Перевод id к типу int
	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id is of invalid type")
		return 0, errors.New("user id is of invalid type")
	}

	// Возвращает полученный id и nil
	return idInt, nil
}
