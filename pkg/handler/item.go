package handler

import (
	"net/http"
	"strconv"
	"todo-app"

	"github.com/gin-gonic/gin"
)

// Хендлер для эндпоинта создания пункта списка
func (h *Handler) createItem(c *gin.Context) {

	// Получение айди пользователя из контекста
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	// Получение айди списка из контекста
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	// Получение вводных данных
	var input todo.TodoItem
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Создание пункта списка
	id, err := h.services.TodoItem.Create(userId, listId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Запись в респонс статус кода 200 и айди нового пункта
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// Хендлер для эндпоинта вывода всех пунктов списка
func (h *Handler) getAllItems(c *gin.Context) {

	// Получение адйди пользователя из контекста
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	// Получение ади списка из контекста
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	// Получение всех пунктов списка
	items, err := h.services.TodoItem.GetAll(userId, listId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Запись в респонс статус кода 200 и всех пунктов
	c.JSON(http.StatusOK, items)
}

// Хендлер для эндпоинта вывода пункта списка по его id
func (h *Handler) getItemById(c *gin.Context) {

	// Получение айди пользвателя из контекста
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	// Получение айди пункта из контекста
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	// Поиск пункта по айди
	item, err := h.services.TodoItem.GetById(userId, itemId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Запись в респонс статус кода 200 и пункта
	c.JSON(http.StatusOK, item)
}

// Хендлер для эндпоинта обновления пункта списка
func (h *Handler) updateItem(c *gin.Context) {

	// Получение айди пользователя из контекста
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	// Получение айди пункта из контекста
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	// Получение новых данных для обновлнения
	var input todo.UpdateItemInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Обновление данных пункта
	if err := h.services.TodoItem.Update(userId, id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Запись в респонс статус кода 200 и "ok"
	c.JSON(http.StatusOK, statusResponse{"ok"})
}

// Хендлер для эндпоинта удаления пункта списка
func (h *Handler) deleteItem(c *gin.Context) {

	// Получение айди пользователя из контекста
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	// Получение айди пункта из контекста
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	// Удаление пункта
	err = h.services.TodoItem.Delete(userId, itemId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Запись в респонс статус кода 200 и "ok"
	c.JSON(http.StatusOK, statusResponse{"ok"})
}
