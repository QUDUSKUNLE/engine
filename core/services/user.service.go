package services

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/db"
	"github.com/medicue/core/domain"
	"github.com/medicue/core/utils"
)

func (service *ServicesHandler) Create(context echo.Context) error {
	dto, ok := context.Get("validatedDTO").(*domain.UserRegisterDTO)
	if !ok {
		return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.InvalidRequestBody), context)
	}
	newUser, err := domain.BuildNewUser(*dto)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}

	user := db.CreateUserParams{
		Email:    newUser.Email,
		Password: newUser.Password,
		UserType: newUser.UserType,
	}
	return service.createUserHelper(
		context, user, db.UserEnumDIAGNOSTICCENTREOWNER, db.UserEnumUSER)
}

func (service *ServicesHandler) CreateDiagnosticCentreManager(context echo.Context) error {
	// Check for permission to add a diagnostic manager
	_, err := utils.PrivateMiddlewareContext(context, string(db.UserEnumDIAGNOSTICCENTREOWNER))
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, context)
	}

	dto, ok := context.Get("validatedDTO").(*domain.DiagnosticCentreManagerRegisterDTO)
	// Check if there are appropriate UserEnumDiagnosticCentreManager
	if !ok || dto.UserType != db.UserEnumDIAGNOSTICCENTREMANAGER {
		return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.InvalidRequestBody), context)
	}
	// Auto generate a password for the manager
	userDto := domain.UserRegisterDTO{
		Email:    dto.Email,
		Password: utils.GenerateRandomPassword(12),
		UserType: dto.UserType,
	}
	newUser, err := domain.BuildNewUser(userDto)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	user := db.CreateUserParams{
		Email:    newUser.Email,
		Password: newUser.Password,
		UserType: newUser.UserType,
	}
	// Create a diagnostic centre manager
	return service.createUserHelper(context, user, db.UserEnumDIAGNOSTICCENTREMANAGER)
}

func (service *ServicesHandler) Login(context echo.Context) error {
	ctx := context.Request().Context()
	dto := context.Get("validatedDTO").(*domain.UserSignInDTO)
	response, err := service.repositoryService.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	if err := domain.ComparePassword(*response, dto.Password); err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	parsedUUID, err := uuid.Parse(response.ID)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	token, err := utils.GenerateToken(domain.CurrentUserDTO{
		UserID: parsedUUID, UserType: response.UserType})
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, context)
	}
	return utils.ResponseMessage(http.StatusOK, map[string]string{"token": token}, context)
}

func (service *ServicesHandler) createUserHelper(
	ctx echo.Context,
	dto db.CreateUserParams,
	allowedTypes ...db.UserEnum,
) error {
	for _, t := range allowedTypes {
		if dto.UserType == t {
			response, err := service.repositoryService.CreateUser(ctx.Request().Context(), dto)
			if err != nil {
				return utils.ErrorResponse(http.StatusBadRequest, err, ctx)
			}
			return utils.ResponseMessage(http.StatusCreated, response, ctx)
		}
	}
	return utils.ErrorResponse(http.StatusBadRequest, errors.New(utils.ErrBadRequest), ctx)
}
