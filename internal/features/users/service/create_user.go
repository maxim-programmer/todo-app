package users_service

import (
	"context"
	"fmt"

	"github.com/nilchan-social/golang-todoapp/internal/core/domain"
)

// CreateUser создаёт нового пользователя: формирует доменный объект,
// валидирует его инварианты и сохраняет через репозиторий.
func (s *UsersService) CreateUser(
	ctx context.Context,
	fullName string,
	phoneNumber *string,
) (domain.User, error) {
	user := domain.CreateUser(
		fullName,
		phoneNumber,
	)

	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}

	user, err := s.usersRepository.SaveUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("save user in repository: %w", err)
	}

	return user, nil
}
