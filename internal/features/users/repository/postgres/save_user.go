package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/nilchan-social/golang-todoapp/internal/core/domain"
)

// SaveUser вставляет нового пользователя в БД и возвращает сохранённую версию.
// RETURNING позволяет получить итоговое состояние записи одним запросом.
func (r *UsersRepository) SaveUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO todoapp.users (id, version, full_name, phone_number)
	VALUES ($1, $2, $3, $4)
	RETURNING id, version, full_name, phone_number;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.ID,
		user.Version,
		user.FullName,
		user.PhoneNumber,
	)

	var userModel UserModel
	if err := userModel.Scan(row); err != nil {
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := modelToDomain(userModel)

	return userDomain, nil
}
