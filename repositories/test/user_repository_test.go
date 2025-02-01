package repositories_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // Import driver PostgreSQL
	"github.com/nabilwafi/warehouse-management-system/models/domain"
	"github.com/nabilwafi/warehouse-management-system/models/web"
	"github.com/nabilwafi/warehouse-management-system/repositories"
	"github.com/stretchr/testify/assert"
)

func TestSaveUser(t *testing.T) {
	db := setupDatabase(t)

	repo := repositories.NewUserRepository()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	ctx := context.Background()

	user := domain.User{
		Email:    "test@example.com",
		Name:     "Test User",
		Role:     "admin",
		Password: "password123",
	}

	err = repo.Save(ctx, tx, &user)
	assert.NoError(t, err, "Save should not return an error")

	var id uuid.UUID
	var email, name, role, password string
	err = tx.QueryRowContext(ctx, "SELECT id, email, name, role, password FROM users WHERE email = $1", "test@example.com").Scan(&id, &email, &name, &role, &password)
	assert.NoError(t, err, "Failed to query user")
	assert.Equal(t, "test@example.com", email, "User email should match")
	assert.Equal(t, "Test User", name, "User name should match")
	assert.Equal(t, "admin", role, "User role should match")
	assert.Equal(t, "password123", password, "User password should match")
}

func TestFindUserByID(t *testing.T) {
	db := setupDatabase(t)

	repo := repositories.NewUserRepository()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	ctx := context.Background()

	userID := uuid.New()
	_, err = tx.ExecContext(ctx, "INSERT INTO users (id, email, name, role, password) VALUES ($1, $2, $3, $4, $5)", userID, "test@example.com", "Test User", "admin", "password123")
	assert.NoError(t, err, "Failed to insert user")

	user, err := repo.FindByID(ctx, tx, userID)
	assert.NoError(t, err, "FindByID should not return an error")
	assert.Equal(t, userID, user.ID, "User ID should match")
	assert.Equal(t, "test@example.com", user.Email, "User email should match")
	assert.Equal(t, "Test User", user.Name, "User name should match")
	assert.Equal(t, domain.UserRole("admin"), user.Role, "User role should match")
}

func TestFindUserByEmail(t *testing.T) {
	db := setupDatabase(t)

	repo := repositories.NewUserRepository()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	ctx := context.Background()

	userID := uuid.New()
	_, err = tx.ExecContext(ctx, "INSERT INTO users (id, email, name, role, password) VALUES ($1, $2, $3, $4, $5)", userID, "test@example.com", "Test User", "admin", "password123")
	assert.NoError(t, err, "Failed to insert user")

	user, err := repo.FindByEmail(ctx, tx, "test@example.com")
	assert.NoError(t, err, "FindByEmail should not return an error")
	assert.Equal(t, userID, user.ID, "User ID should match")
	assert.Equal(t, "test@example.com", user.Email, "User email should match")
	assert.Equal(t, "Test User", user.Name, "User name should match")
	assert.Equal(t, domain.UserRole("admin"), user.Role, "User role should match")
	assert.Equal(t, "password123", user.Password, "User password should match")
}

func TestFindAllUsers(t *testing.T) {
	db := setupDatabase(t)

	repo := repositories.NewUserRepository()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	ctx := context.Background()

	// Insert beberapa data user untuk testing
	_, err = tx.ExecContext(ctx, "INSERT INTO users (id, email, name, role, password) VALUES ($1, $2, $3, $4, $5)", uuid.New(), "test1@example.com", "Test User 1", "admin", "password123")
	assert.NoError(t, err, "Failed to insert user")

	_, err = tx.ExecContext(ctx, "INSERT INTO users (id, email, name, role, password) VALUES ($1, $2, $3, $4, $5)", uuid.New(), "test2@example.com", "Test User 2", "staff", "password456")
	assert.NoError(t, err, "Failed to insert user")

	userQuery := &web.UserQuery{}
	users, err := repo.FindAll(ctx, tx, userQuery)
	assert.NoError(t, err, "FindAll should not return an error")
	assert.Equal(t, 2, len(users), "Should return 2 users")

	// Verifikasi data yang diambil
	assert.Equal(t, "test1@example.com", users[0].Email, "First user email should match")
	assert.Equal(t, "Test User 1", users[0].Name, "First user name should match")
	assert.Equal(t, domain.UserRole("admin"), users[0].Role, "First user role should match")
	assert.Equal(t, "password123", users[0].Password, "First user password should match")

	assert.Equal(t, "test2@example.com", users[1].Email, "Second user email should match")
	assert.Equal(t, "Test User 2", users[1].Name, "Second user name should match")
	assert.Equal(t, domain.UserRole("staff"), users[1].Role, "Second user role should match")
	assert.Equal(t, "password456", users[1].Password, "Second user password should match")
}
