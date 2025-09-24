package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"speakpall/api/models"
	"speakpall/pkg/logger"
	"speakpall/storage"
)

type settingsRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewSettingsRepo(db *pgxpool.Pool, log logger.ILogger) storage.ISettingsStorage {
	return &settingsRepo{db: db, log: log}
}

func (r *settingsRepo) GetUserSettings(ctx context.Context, userID string) (*models.UserSettings, error) {
	const q = `
SELECT discoverable, allow_messages, notify_push, notify_email,
       COALESCE(to_char(updated_at,'YYYY-MM-DD"T"HH24:MI:SS"Z"'),'') AS updated_at
FROM user_settings WHERE user_id=$1`
	var s models.UserSettings
	err := r.db.QueryRow(ctx, q, userID).Scan(
		&s.Discoverable, &s.AllowMessages, &s.NotifyPush, &s.NotifyEmail, &s.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			// default settings
			return &models.UserSettings{
				Discoverable:  true,
				AllowMessages: true,
				NotifyPush:    true,
				NotifyEmail:   false,
				UpdatedAt:     "",
			}, nil
		}
		r.log.Error("GetUserSettings: query failed", logger.Error(err), logger.String("user_id", userID))
		return nil, err
	}
	return &s, nil
}

func (r *settingsRepo) UpsertUserSettings(ctx context.Context, userID string, req models.UpdateSettingsRequest) error {
	// old values (or defaults)
	cur, _ := r.GetUserSettings(ctx, userID)

	if req.Discoverable != nil {
		cur.Discoverable = *req.Discoverable
	}
	if req.AllowMessages != nil {
		cur.AllowMessages = *req.AllowMessages
	}
	if req.NotifyPush != nil {
		cur.NotifyPush = *req.NotifyPush
	}
	if req.NotifyEmail != nil {
		cur.NotifyEmail = *req.NotifyEmail
	}

	const q = `
INSERT INTO user_settings(user_id, discoverable, allow_messages, notify_push, notify_email, updated_at)
VALUES ($1,$2,$3,$4,$5, now())
ON CONFLICT (user_id) DO UPDATE SET
  discoverable = EXCLUDED.discoverable,
  allow_messages = EXCLUDED.allow_messages,
  notify_push = EXCLUDED.notify_push,
  notify_email = EXCLUDED.notify_email,
  updated_at = now()`
	_, err := r.db.Exec(ctx, q, userID, cur.Discoverable, cur.AllowMessages, cur.NotifyPush, cur.NotifyEmail)
	if err != nil {
		r.log.Error("UpsertUserSettings: exec failed", logger.Error(err), logger.String("user_id", userID))
		return err
	}
	return nil
}
