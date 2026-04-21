// Package users_service содержит бизнес-логику управления пользователями.
package users_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nilchan-social/golang-todoapp/internal/core/domain"
)

// UsersService — сервис пользователей с бизнес-логикой CRUD-операций.
type UsersService struct {
	usersRepository UsersRepository
}

// UsersRepository — интерфейс репозитория пользователей.
// Определён в пакете сервиса по принципу Dependency Inversion.
type UsersRepository interface {
	SaveUser(
		ctx context.Context,
		user domain.User,
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

	UpdateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
}

// NewUsersService создаёт сервис пользователей с внедрённым репозиторием.
func NewUsersService(
	usersRepository UsersRepository,
) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}
