package handler

// хендлеры, которые отвечают за авторизацию, - привязываються к основному классу хендлеров по ссылке
import (
	"net/http"
	"todo"

	"github.com/gin-gonic/gin"
)


func (h *Handler) signUp (c *gin.Context) {
	var input todo.User

	// получаем данные из тела запроса, десериализуем и привязываем их к обьекту инпут по указателю
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// обращаемся к сервису авторизации и создаем пользователя
	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// отправляем ответ на фронтенд
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}


type signInInput struct {
	// валидация- обязательное поле в теле запроса по флагам binding, и наименование поля в джейсоне для сериализации и десериализации
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn (c *gin.Context) {
	var input signInInput

		// получаем данные из тела запроса, десериализуем и привязываем их к обьекту инпут по указателю
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// гененириуем токен через сервис просле логина
	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

		// отправляем ответ на фронтенд
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}