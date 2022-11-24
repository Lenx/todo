package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lenx/todo"
)

func (h *Handler) signUp(c *gin.Context) {
	var input todo.User

	// получаем JSON и упаковываем в структуру
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	// запускаем метод CreateUser у сервиса, получаем id созданного юзера
	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// подготавливаем JSON с ответом
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	// получаем JSON и упаковываем в структуру
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// генерируем JWT-токен
	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// подготавливаем JSON с ответом
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
