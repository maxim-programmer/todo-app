package users_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nilchan-social/golang-todoapp/internal/core/domain"
)

// PatchUser частично обновляет пользователя по ID.
// Паттерн «read-modify-write»: читаем → применяем патч → сохраняем.
// Оптимистичная блокировка на уровне репозитория защищает от конкурентных изменений.
func (s *UsersService) PatchUser(
	ctx context.Context,
	id uuid.UUID,
	patch domain.UserPatch,
) (domain.User, error) {
	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user from repository: %w", err)
	}

	if err := user.ApplyPatch(patch); err != nil {
		return domain.User{}, fmt.Errorf("apply user patch: %w", err)
	}

	patchedUser, err := s.usersRepository.UpdateUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("update user in repository: %w", err)
	}

	return patchedUser, nil
}
