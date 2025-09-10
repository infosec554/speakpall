package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"speakpall/api/models"
	"speakpall/pkg/logger"
	"speakpall/storage"
)

type userRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewUserRepo(db *pgxpool.Pool, log logger.ILogger) storage.IUserStorage {
	return &userRepo{
		db:  db,
		log: log,
	}
}
func (r *userRepo) Create(ctx context.Context, req models.SignupRequest) (string, error) {
	id := uuid.New().String()
	query := `
		INSERT INTO users (id, name, email, password, status, role)
		VALUES ($1, $2, $3, $4, 'active', 'user')
	`
	_, err := r.db.Exec(ctx, query, id, req.Name, req.Email, req.Password)
	if err != nil {
		r.log.Error("error inserting user", logger.Error(err))
		return "", err
	}
	return id, nil
}
func (r *userRepo) GetForLoginByEmail(ctx context.Context, email string) (models.LoginUser, error) {
	var user models.LoginUser
	query := `
		SELECT id, password, status, role
		FROM users
		WHERE email = $1 AND status = 'active'
	`
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Password,
		&user.Status,
		&user.Role,
	)
	if err != nil {
		r.log.Error("failed to get user by email", logger.Error(err))
		return models.LoginUser{}, err
	}
	return user, nil
}
func (r *userRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, name, email, status, role, created_at
		FROM users
		WHERE id = $1 AND status = 'active'
	`
	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Status,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		r.log.Error("failed to get user by ID", logger.Error(err))
		return nil, err
	}
	return &user, nil
}

// postgres/userRepo.go

func (r *userRepo) UpdatePassword(ctx context.Context, userID, newPassword string) error {
	query := `UPDATE users SET password = $1 WHERE id = $2 AND status = 'active'`
	_, err := r.db.Exec(ctx, query, newPassword, userID)
	if err != nil {
		r.log.Error("failed to update user password", logger.Error(err))
		return err
	}
	return nil
}

func (r *userRepo) GetPasswordByID(ctx context.Context, userID string) (string, error) {
	var hashedPassword string
	query := `SELECT password FROM users WHERE id = $1 AND status = 'active'`
	err := r.db.QueryRow(ctx, query, userID).Scan(&hashedPassword)
	if err != nil {
		r.log.Error("failed to get user password", logger.Error(err))
		return "", err
	}
	return hashedPassword, nil
}
