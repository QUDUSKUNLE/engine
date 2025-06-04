package services

import (
	// "go/token"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/db"
	"github.com/medicue/core/domain"
	"github.com/medicue/core/utils"
)

func (service *ServicesHandler) Create(context echo.Context) error {
	ctx := context.Request().Context()
	dto := context.Get("validatedDTO").(domain.UserRegisterDTO)
	newUser, err := domain.BuildNewUser(dto)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	user := db.CreateUserParams{
		Email:    newUser.Email,
		Password: newUser.Password,
		UserType: newUser.UserType,
	}
	response, err := service.repositoryService.CreateUser(ctx, user)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	return utils.ResponseMessage(http.StatusCreated, response, context)
}

func (service *ServicesHandler) Login(context echo.Context) error {
	ctx := context.Request().Context()
	dto := context.Get("validatedDTO").(domain.UserSignInDTO)
	response, err := service.repositoryService.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	if err := domain.ComparePassword(*response, dto.Password); err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	token, err := utils.GenerateToken(domain.CurrentUserDTO{UserID: response.ID, UserType: string(response.UserType)})
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	return utils.ResponseMessage(http.StatusCreated, map[string]string{"token": token}, context)
}
