// storage/postgres/user_repo.go
package postgres

import (
	"context"
	"time"

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
	return &userRepo{db: db, log: log}
}

// CreateUser: req.Password bu yerga KELGUNCHA hashlangan bo‘lishi kerak.
func (r *userRepo) CreateUser(ctx context.Context, req models.SignupRequest) (string, error) {
	id := uuid.New().String()
	const q = `
		INSERT INTO users (
			id, email, display_name, password_hash, target_lang, level, country_code
		) VALUES ($1,$2,$3,$4,$5,$6,$7)
	`
	_, err := r.db.Exec(ctx, q,
		id,
		req.Email,
		req.DisplayName,
		req.Password, // hashlangan parol
		req.TargetLang,
		req.Level,
		req.CountryCode,
	)
	if err != nil {
		r.log.Error("users insert failed", logger.Error(err))
		return "", err
	}
	return id, nil
}

func (r *userRepo) GetLoginByEmail(ctx context.Context, email string) (models.LoginUser, error) {
	const q = `SELECT id, password_hash, role FROM users WHERE email = $1`
	var u models.LoginUser
	if err := r.db.QueryRow(ctx, q, email).Scan(&u.ID, &u.PasswordHash, &u.Role); err != nil {
		r.log.Error("get login by email failed", logger.Error(err))
		return models.LoginUser{}, err
	}
	return u, nil
}

func (r *userRepo) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	const q = `
		SELECT
			id, email, display_name, password_hash, google_id, avatar_url,
			age, gender, country_code, target_lang, level, role,
			created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var u models.User
	if err := r.db.QueryRow(ctx, q, id).Scan(
		&u.ID, &u.Email, &u.DisplayName, &u.PasswordHash, &u.GoogleID, &u.AvatarURL,
		&u.Age, &u.Gender, &u.CountryCode, &u.TargetLang, &u.Level, &u.Role,
		&u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		r.log.Error("get user by id failed", logger.Error(err))
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) UpdatePasswordHash(ctx context.Context, userID, newHash string) error {
	const q = `UPDATE users SET password_hash=$1, updated_at=NOW() WHERE id=$2`
	if _, err := r.db.Exec(ctx, q, newHash, userID); err != nil {
		r.log.Error("update password failed", logger.Error(err))
		return err
	}
	return nil
}

func (r *userRepo) UpdateRole(ctx context.Context, userID, role string) error {
	const q = `UPDATE users SET role=$1, updated_at=NOW() WHERE id=$2`
	if _, err := r.db.Exec(ctx, q, role, userID); err != nil {
		r.log.Error("update role failed", logger.Error(err))
		return err
	}
	return nil
}

// --- Password reset (repo token yaratmaydi) ---

func (r *userRepo) SavePasswordResetToken(ctx context.Context, userID, token string, expiresAt time.Time) error {
	const q = `
		INSERT INTO password_reset_tokens (user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, NOW())
	`
	if _, err := r.db.Exec(ctx, q, userID, token, expiresAt); err != nil {
		r.log.Error("save reset token failed", logger.Error(err))
		return err
	}
	return nil
}

func (r *userRepo) GetUserIDByPasswordResetToken(ctx context.Context, token string, now time.Time) (string, error) {
	const q = `
		SELECT user_id
		FROM password_reset_tokens
		WHERE token = $1 AND expires_at > $2
		ORDER BY created_at DESC
		LIMIT 1
	`
	var userID string
	if err := r.db.QueryRow(ctx, q, token, now).Scan(&userID); err != nil {
		r.log.Error("validate reset token failed", logger.Error(err))
		return "", err
	}
	return userID, nil
}

func (r *userRepo) GetPasswordByID(ctx context.Context, userID string) (string, error) {
	const q = `SELECT password_hash FROM users WHERE id = $1`
	var ph string
	if err := r.db.QueryRow(ctx, q, userID).Scan(&ph); err != nil {
		r.log.Error("get password_hash by id failed", logger.Error(err))
		return "", err
	}
	return ph, nil
}
