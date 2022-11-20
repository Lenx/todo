package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	todo "github.com/lenx/todo/pkg"
)

func (h *Handler) signUp(c *gin.Context) {

	/*
			var input todo.User

		    if err := c.ShouldBind(&input); err != nil {
		        c.AbortWithError(http.StatusBadRequest, err)
		        return
		    }
	*/

	var input todo.User

	// ловим JSON и упаковываем в структуру
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

func (h *Handler) signUpShow(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", nil)
}

func (h *Handler) signInShow(c *gin.Context) {
	c.HTML(http.StatusOK, "signin.html", nil)
}

type signInInput struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password"json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {

	var input signInInput

	// ловим JSON и упаковываем в структуру
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
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
