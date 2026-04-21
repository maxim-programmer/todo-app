package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nilchan-social/golang-todoapp/internal/core/domain"
	core_errors "github.com/nilchan-social/golang-todoapp/internal/core/errors"
	core_postgres_pool "github.com/nilchan-social/golang-todoapp/internal/core/repository/postgres/pool"
)

// GetUser возвращает пользователя по ID.
// ErrNoRows транслируется в ErrNotFound → HTTP 404.
func (r *UsersRepository) GetUser(
	ctx context.Context,
	id uuid.UUID,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, full_name, phone_number
	FROM todoapp.users
	WHERE id=$1;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var userModel UserModel
	if err := userModel.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id='%s': %w",
				id,
				core_errors.ErrNotFound,
			)
		}

		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := modelToDomain(userModel)

	return userDomain, nil
}
