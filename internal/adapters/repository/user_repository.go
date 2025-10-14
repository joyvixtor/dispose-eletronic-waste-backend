package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/joyvixtor/dispose-eletronic-waste-backend/internal/domain/models"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	// Gera UUID se não fornecido
	if user.Id == "" {
		user.Id = uuid.NewString()
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := `
		INSERT INTO users (id, full_name, ufpe_email, password, workplace, role)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.ExecContext(ctx, query,
		user.Id,
		user.FullName,
		user.UFPEEmail,
		string(hashedPassword),
		user.Workplace,
		user.Role,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil

}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, full_name, ufpe_email, password, workplace, role
		FROM users
		WHERE ufpe_email = ?
	`

	row := r.db.QueryRowContext(ctx, query, email)

	user := &models.User{}
	err := row.Scan(&user.Id, &user.FullName, &user.UFPEEmail, &user.Password, &user.Workplace, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (r *UserRepository) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
