package repositories

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
)

type UserRepository interface {
	Save(register *dto.RegisterRequest) (err error)
	FindByEmail(email string) (user *dto.User, code int, err error)
	FindByID(id uuid.UUID) (user *dto.User, code int, err error)
	GetAll(pagination *web.PaginationRequest) ([]*dto.User, error)
}

type userRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

func (r *userRepositoryImpl) GetAll(pagination *web.PaginationRequest) ([]*dto.User, error) {
	var userData []*dto.User

	offset := (pagination.Page - 1) * pagination.Size

	rows, err := r.db.Queryx("SELECT id, email, name, role FROM public.users OFFSET $1 LIMIT $2", offset, pagination.Size)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user dto.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Role); err != nil {
			return nil, err
		}
		userData = append(userData, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userData, nil
}

func (r *userRepositoryImpl) Save(register *dto.RegisterRequest) (err error) {
	uuid := uuid.New()

	_, err = r.db.Exec("INSERT INTO public.users (id, email, password, name, role) VALUES ($1, $2, $3, $4, $5)", uuid, register.Email, register.Password, register.Name, register.Role)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepositoryImpl) FindByEmail(email string) (user *dto.User, code int, err error) {
	var userData dto.User

	if err := r.db.QueryRow("SELECT id, email, name, role, password FROM public.users WHERE email = $1", email).Scan(&userData.ID, &userData.Email, &userData.Name, &userData.Role, &userData.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, 404, sql.ErrNoRows
		}

		return nil, 500, err
	}

	return &userData, 200, nil
}

func (r *userRepositoryImpl) FindByID(id uuid.UUID) (user *dto.User, code int, err error) {
	var userData dto.User

	if err := r.db.QueryRow("SELECT id, email, name, role, password FROM public.users WHERE id = $1", id).Scan(&userData.ID, &userData.Email, &userData.Name, &userData.Role, &userData.Password); err != nil {
		if sql.ErrNoRows != nil {
			return nil, 404, sql.ErrNoRows
		}

		return nil, 500, err
	}

	return &userData, 200, nil
}
