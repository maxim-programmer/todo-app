package users_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nilchan-social/golang-todoapp/internal/core/domain"
)

// GetUser возвращает пользователя по ID, делегируя запрос репозиторию.
func (s *UsersService) GetUser(
	ctx context.Context,
	id uuid.UUID,
) (domain.User, error) {
	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user from repository: %w", err)
	}

	return user, nil
}
