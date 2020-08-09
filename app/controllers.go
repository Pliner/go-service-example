package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"net/http"
)

type UserModel struct {
	Id        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

type UsersController struct {
	usersRepository *UsersRepository
	logger          *zap.Logger
}

func (controller *UsersController) SaveUser(context *gin.Context) {
	user := &UserModel{}
	if err := context.ShouldBindJSON(&user); err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	err := controller.usersRepository.SaveUser(context.Request.Context(), user.Id, user.FirstName, user.LastName)
	if err != nil {
		context.String(http.StatusInternalServerError, err.Error())
		return
	}
	context.Status(http.StatusOK)
}

func (controller *UsersController) SelectUsers(context *gin.Context) {
	users, err := controller.usersRepository.SelectUsers(context.Request.Context())
	if err != nil {
		context.String(http.StatusInternalServerError, err.Error())
		return
	}
	context.JSON(http.StatusOK, convertFromUserRows(users))
}

func (controller *UsersController) GetUserById(context *gin.Context) {
	userId, err := uuid.FromString(context.Param("id"))
	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	user, err := controller.usersRepository.GetUserById(context.Request.Context(), userId)
	if err != nil {
		context.String(http.StatusInternalServerError, err.Error())
		return
	}
	context.JSON(http.StatusOK, convertFromUserRow(user))
}

func (controller *UsersController) DeleteUserById(context *gin.Context) {
	userId, err := uuid.FromString(context.Param("id"))
	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	err = controller.usersRepository.DeleteUserById(context.Request.Context(), userId)
	if err != nil {
		context.String(http.StatusInternalServerError, err.Error())
		return
	}
	context.JSON(http.StatusOK, nil)
}

func convertFromUserRow(user *User) *UserModel {
	return &UserModel{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}

func convertFromUserRows(users []*User) []*UserModel {
	result := make([]*UserModel, 0)
	for _, user := range users {
		result = append(result, convertFromUserRow(user))
	}
	return result
}
