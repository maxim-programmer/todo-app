// Package users_transport_http содержит HTTP-обработчики для фичи пользователей.
package users_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/nilchan-social/golang-todoapp/internal/core/domain"
	core_http_server "github.com/nilchan-social/golang-todoapp/internal/core/transport/http/server"
)

// UsersHTTPHandler — HTTP-обработчик для операций с пользователями.
type UsersHTTPHandler struct {
	usersService UsersService
}

// UsersService — интерфейс сервиса пользователей.
// Определён в транспортном слое по принципу Dependency Inversion.
type UsersService interface {
	CreateUser(
		ctx context.Context,
		fullName string,
		phoneNumber *string,
	) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)

	GetUser(
		ctx context.Context,
		id uuid.UUID,
	) (domain.User, error)

	DeleteUser(
		ctx context.Context,
		id uuid.UUID,
	) error

	PatchUser(
		ctx context.Context,
		id uuid.UUID,
		patch domain.UserPatch,
	) (domain.User, error)
}

// NewUsersHTTPHandler создаёт обработчик пользователей с внедрённым сервисом.
func NewUsersHTTPHandler(
	usersService UsersService,
) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}

// Routes возвращает маршруты REST API для пользователей.
func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: h.GetUsers,
			/*
				Example of usage Middleware on separate Route

				Middleware: []core_http_middleware.Middleware{
					core_http_middleware.Dummy("get users middleware"),
				},
			*/
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{id}",
			Handler: h.GetUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/{id}",
			Handler: h.DeleteUser,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/{id}",
			Handler: h.PatchUser,
		},
	}
}
