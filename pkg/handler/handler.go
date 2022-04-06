package handler

import (
	"todo-app/pkg/service"

	"github.com/gin-gonic/gin"
)

// Структура хендлера
type Handler struct {
	services *service.Service
}

// Конструктор хендлера
func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

// Функция инициализации эндпоинтов
func (h *Handler) InitRoutes() *gin.Engine {

	// Инициализация роутера
	router := gin.New()

	// Группа эндпоинтов для аутентификации
	auth := router.Group("/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)
	}

	// Группа эндпоинтов для работы со списками
	api := router.Group("/api", h.userIdentity)
	{
		lists := api.Group("/lists")
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)
			lists.GET("/:id", h.getListById)
			lists.PUT("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)

			items := lists.Group(":id/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItems)
			}
		}
	}

	// Группа эндпоинтов для работы с пунктами списка
	items := api.Group("items")
	{
		items.GET("/:id", h.getItemById)
		items.PUT("/:id", h.updateItem)
		items.DELETE("/:id", h.deleteItem)
	}

	// Возвращает готовый роутер
	return router
}
