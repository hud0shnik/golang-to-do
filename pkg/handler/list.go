package handler

import (
	"net/http"
	"strconv"
	"todo-app"

	"github.com/gin-gonic/gin"
)

// Хендлер для эндпоинта создания списка
func (h *Handler) createList(c *gin.Context) {

	// Получение id пользователя из контекста
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Данные для создания списка
	var input todo.TodoList
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Создание списка
	id, err := h.services.TodoList.Create(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Запись в ответ статус кода 200 и id списка
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// Структура для вывода всех списков
type getAllListsResponse struct {
	Data []todo.TodoList `json:"data"`
}

// Хендлер для эндпоинта вывода всех списков пользователя
func (h *Handler) getAllLists(c *gin.Context) {

	// Получение id пользователя из контекста
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Получение всех списков пользователя
	lists, err := h.services.TodoList.GetAll(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Запись в ответ статус кода 200 и найденных списков
	c.JSON(http.StatusOK, getAllListsResponse{
		Data: lists,
	})

}

// Хендлер для эндпоинта получения списка по id
func (h *Handler) getListById(c *gin.Context) {

	// Получение id пользователя из контекста
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//  Получение id списка из пути запроса
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	// Получение списка
	list, err := h.services.TodoList.GetById(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Запись в ответ статус кода 200 и списка
	c.JSON(http.StatusOK, list)
}

// Хендлер для эндпоинта обновления списка
func (h *Handler) updateList(c *gin.Context) {

	// Получение id пользователя из контекста
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Получение id списка из параметра запроса
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	// Получение вводных данных для обновленного списка
	var input todo.UpdateListInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Обновление списка
	if err := h.services.TodoList.Update(userId, id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Запись в ответ статус кода 200 и статуса ок
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// Хендлер для эндпоинта удаления списка
func (h *Handler) deleteList(c *gin.Context) {

	// Получение id пользователя из контекста
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Получение id списка из параметра запроса
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	// Удаление списка
	err = h.services.TodoList.Delete(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Запись в ответ статус кода 200 и статуса ок
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
