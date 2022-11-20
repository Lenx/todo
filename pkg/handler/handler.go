package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lenx/todo/pkg/service"
)

// структура, необходимая для установления зависимостей
type Handler struct {
	services *service.Service
}

// создаем объект Handler
func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

// создаем роутер(gin) и инициализируем эндпоинты
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Static("/css", "./templates/css")
	router.LoadHTMLGlob("templates/*.html")

	start := router.Group("/")
	{
		start.GET("/", h.listsShow)
	}

	auth := router.Group("/auth")
	{
		auth.GET("/sign-up", h.signUpShow)
		auth.GET("/sign-in", h.signInShow)
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	lists := router.Group("/user_lists")
	{
		lists.GET("/", h.listsShow)
	}

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

		items := api.Group("/items")
		{
			items.GET("/:id", h.getItemById)
			items.PUT("/:id", h.updateItem)
			items.DELETE("/:id", h.deleteItem)
			tags := items.Group(":id/tags")
			{
				tags.POST("/", h.createTag)
				tags.GET("/", h.getAllTags)
			}
		}
	}

	return router
}
