package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	userusecase "github.com/hkobori/golang-domain-driven-arch/internal/app/usecase/user"
)

type UserHandler struct {
	createUser userusecase.CreateUserUseCase
}

func NewUserHandler(createUser userusecase.CreateUserUseCase) *UserHandler {
	return &UserHandler{createUser: createUser}
}

type createUserRequest struct {
	Name  string `json:"name"  validate:"required"`
	Email string `json:"email" validate:"required,email"`
}


type userResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func toResponse(out *userusecase.UserOutput) *userResponse {
	return &userResponse{
		ID:    out.ID,
		Name:  out.Name,
		Email: out.Email,
	}
}

func (h *UserHandler) Create(c echo.Context) error {
	var req createUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	out, err := h.createUser.Execute(c.Request().Context(), userusecase.CreateUserInput{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		var ucErr *userusecase.CreateUserError
		if errors.As(err, &ucErr) {
			switch ucErr.Kind {
			case userusecase.ErrValidation:
				return echo.NewHTTPError(http.StatusBadRequest, ucErr.Message)
			case userusecase.ErrEmailDuplicated:
				return echo.NewHTTPError(http.StatusConflict, ucErr.Message)
			}
		}
		return err
	}

	return c.JSON(http.StatusCreated, toResponse(out))
}
