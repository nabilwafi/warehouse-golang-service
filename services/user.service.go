package services

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nabilwafi/warehouse-management-system/exception"
	"github.com/nabilwafi/warehouse-management-system/models/domain"
	"github.com/nabilwafi/warehouse-management-system/models/web"
	"github.com/nabilwafi/warehouse-management-system/repositories"
	"github.com/nabilwafi/warehouse-management-system/utils"
)

type UserServiceImpl struct {
	repo      repositories.UserRepository
	validator *validator.Validate
	db        *sql.DB
}

func NewUserService(repo repositories.UserRepository, validator *validator.Validate, db *sql.DB) UserService {
	return &UserServiceImpl{
		repo:      repo,
		validator: validator,
		db:        db,
	}
}

func (svc *UserServiceImpl) Login(ctx context.Context, user *web.Login) (string, error) {
	err := svc.validator.StructCtx(ctx, user)
	if err != nil {
		return "", exception.NewCustomError(400, err.Error())
	}

	tx, err := svc.db.Begin()
	if err != nil {
		return "", exception.NewCustomError(500, "failed to begin transaction: "+err.Error())
	}
	defer tx.Rollback()

	userData, err := svc.repo.FindByEmail(ctx, tx, user.Email)
	if err != nil {
		return "", exception.NewCustomError(404, "user not found")
	}

	if err := utils.ComparePasswords(userData.Password, user.Password); err != nil {
		return "", exception.NewCustomError(401, "invalid email or password")
	}

	token, err := utils.GenerateToken(userData.ID, userData.Name, string(userData.Role))
	if err != nil {
		return "", exception.NewCustomError(500, "invalid email or password")
	}

	if err := tx.Commit(); err != nil {
		return "", exception.NewCustomError(500, "failed to commit transaction: "+err.Error())
	}

	return token, nil
}

func (svc *UserServiceImpl) Register(ctx context.Context, user *web.Register) error {
	err := svc.validator.StructCtx(ctx, user)
	if err != nil {
		return exception.NewCustomError(400, err.Error())
	}

	tx, err := svc.db.Begin()
	if err != nil {
		return exception.NewCustomError(500, "failed to begin transaction: "+err.Error())
	}
	defer tx.Rollback()

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return exception.NewCustomError(500, "failed to hash password: "+err.Error())
	}

	newUser := &domain.User{
		ID:       uuid.New(),
		Name:     user.Name,
		Email:    user.Email,
		Password: hashedPassword,
		Role:     user.Role,
	}

	if err := svc.repo.Save(ctx, tx, newUser); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return exception.NewCustomError(500, "failed to commit transaction: "+err.Error())
	}

	return nil
}

func (svc *UserServiceImpl) GetMe(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	tx, err := svc.db.Begin()
	if err != nil {
		return nil, exception.NewCustomError(500, "failed to begin transaction: "+err.Error())
	}
	defer tx.Rollback()

	userData, err := svc.repo.FindByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, exception.NewCustomError(500, "failed to commit transaction: "+err.Error())
	}

	return userData, nil
}

func (svc *UserServiceImpl) ListUser(ctx context.Context, userQuery *web.UserQuery) ([]domain.User, error) {
	tx, err := svc.db.Begin()
	if err != nil {
		return nil, exception.NewCustomError(500, "failed to begin transaction: "+err.Error())
	}
	defer tx.Rollback()

	users, err := svc.repo.FindAll(ctx, tx, userQuery)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, exception.NewCustomError(500, "failed to commit transaction: "+err.Error())
	}

	return users, nil
}
