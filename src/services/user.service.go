package services

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
	"github.com/nabilwafi/warehouse-management-system/src/repositories"
	"github.com/nabilwafi/warehouse-management-system/src/utils"
)

type UserService interface {
	Login(login *dto.LoginRequest) (string, int, error)
	Register(register *dto.RegisterRequest) (int, error)
	GetUserByID(id uuid.UUID) (*dto.User, int, error)
	GetAllUser(pagination *web.PaginationRequest) ([]*dto.User, int, error)
}

type UserServiceImpl struct {
	user repositories.UserRepository
}

func NewUserService(user repositories.UserRepository) UserService {
	return &UserServiceImpl{
		user: user,
	}
}

func (s *UserServiceImpl) Login(login *dto.LoginRequest) (string, int, error) {
	user, code, err := s.user.FindByEmail(login.Email)
	if err != nil {
		return "", code, err
	}

	err = utils.VerifyPassword(user.Password, login.Password)
	if err != nil {
		return "", 400, errors.New("wrong password")
	}

	claims := &utils.CustomClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	tokenString, err := utils.GenerateToken(claims)
	if err != nil {
		return "", 500, err
	}

	return tokenString, 200, nil
}

func (s *UserServiceImpl) Register(register *dto.RegisterRequest) (int, error) {
	user, code, err := s.user.FindByEmail(register.Email)
	if err != nil {
		if err != sql.ErrNoRows || code == 500 {
			return code, err
		}
	}

	if user != nil {
		return 401, errors.New("email is exists")
	}

	hashedPassword, err := utils.HashPassword(register.Password)
	if err != nil {
		return 500, err
	}

	register.Password = hashedPassword

	err = s.user.Save(register)
	if err != nil {
		return 500, err
	}

	return 201, nil
}

func (s *UserServiceImpl) GetUserByID(id uuid.UUID) (*dto.User, int, error) {
	user, code, err := s.user.FindByID(id)
	if err != nil {
		return nil, code, err
	}

	return user, 200, nil
}

func (s *UserServiceImpl) GetAllUser(pagination *web.PaginationRequest) ([]*dto.User, int, error) {
	users, err := s.user.GetAll(pagination)
	if err != nil {
		return nil, 500, err
	}

	return users, 200, nil
}
