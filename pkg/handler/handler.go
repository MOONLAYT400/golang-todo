// создаем пакет хендлеров
package handler

import (
	"todo/pkg/service"

	"github.com/gin-gonic/gin"
)

// основная структура handler, аналог объекта и класса  в js одновременно
type Handler struct{
	services *service.Service
}

// конструктор handler-а, как конструктор класса, создающий объект, функция-фабрика
func NewHandler(services *service.Service) *Handler{
	return &Handler{services: services} //new handler constructor with service
}

// инициализация маршрутов в хендлере, вызывается в main, синтаксис (h *Handler) показывает, что мы таким образом добавляем метод класса напряму в объект
// по указателю, привязываем метод к конретному объекту 
func (h *Handler) InitRoutes() *gin.Engine {
	// создаем роутер с помощью пакета jin
	router:= gin.New()

	// группы маршрутов
	auth:= router.Group("/auth")
	{
		// адрес маршрута и его обработчик
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}
// на группу api привязываем middleware проверки токена - получаем не публичные эндпоинты
	api:= router.Group("/api",h.userIdentity)
	{
		lists:= api.Group("/lists")
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)
			lists.GET("/:id", h.getListById)
			lists.PUT("/:id",h.updateList)
			lists.DELETE("/:id",h.deleteList)

			items:= lists.Group(":id/items")
			{
				items.POST("/",h.createItem)
				items.GET("/",h.getAllItems)

			}
		}
		
		items := api.Group("/items")
		{
				items.GET("/:item_id",h.getItemById)
				items.PUT("/:item_id",h.updateItem)
				items.DELETE("/:item_id",h.deleteItem)
		}
		}

		return router
}