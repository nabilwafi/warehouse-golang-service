package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/nabilwafi/warehouse-management-system/models/domain"
	"github.com/nabilwafi/warehouse-management-system/models/web"
	"github.com/nabilwafi/warehouse-management-system/utils"
)

type userRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &userRepositoryImpl{}
}

func (r *userRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, user *domain.User) error {
	uuid := utils.GenerateUUID()

	query := `INSERT INTO users(id, email, name, role, password) VALUES ($1,$2,$3,$4,$5)`
	result, err := tx.ExecContext(ctx, query, uuid, user.Email, user.Name, user.Role, user.Password)
	if err != nil {
		log.Println(err.Error())
		return errors.New("failed to insert user: " + err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println(err.Error())
		return errors.New("failed to get affected rows: " + err.Error())
	}

	if rowsAffected == 0 {
		return errors.New("no rows inserted")
	}

	return nil
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, tx *sql.Tx, id uuid.UUID) (*domain.User, error) {
	query := `SELECT id, email, name, role, password FROM users WHERE id = $1`
	row := tx.QueryRowContext(ctx, query, id)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			log.Println(err.Error())
			return nil, errors.New("user not found")
		}

		log.Println(err.Error())
		return nil, errors.New("failed to scan user: " + err.Error())
	}

	return &user, nil
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, tx *sql.Tx, email string) (*domain.User, error) {
	query := `SELECT id, email, name, role, password FROM users WHERE email = $1`
	row := tx.QueryRowContext(ctx, query, email)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			log.Println(err.Error())
			return nil, errors.New("user not found")
		}
		log.Println(err.Error())
		return nil, errors.New("failed to scan user: " + err.Error())
	}

	return &user, nil
}

func (r *userRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, userQuery *web.UserQuery) ([]domain.User, error) {
	query := `SELECT id, email, name, role, password FROM users`
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("failed to query users: " + err.Error())
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.Password); err != nil {
			log.Println(err.Error())
			return nil, errors.New("failed to scan user: " + err.Error())
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Println(err.Error())
		return nil, errors.New("error occurred during rows iteration: " + err.Error())
	}

	return users, nil
}
