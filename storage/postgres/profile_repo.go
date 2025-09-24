package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"speakpall/api/models"
	"speakpall/pkg/logger"
	"speakpall/storage"
)

type profileRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewProfileRepo(db *pgxpool.Pool, log logger.ILogger) storage.IProfileStorage {
	return &profileRepo{db: db, log: log}
}

func (r *profileRepo) GetProfile(ctx context.Context, userID string) (*models.Profile, error) {
	const q = `
SELECT id, email, display_name, avatar_url, age, gender, country_code,
       native_lang, target_lang, level, about, timezone,
       to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS created_at
FROM users
WHERE id = $1`
	var p models.Profile
	err := r.db.QueryRow(ctx, q, userID).Scan(
		&p.ID, &p.Email, &p.DisplayName, &p.AvatarURL, &p.Age, &p.Gender, &p.CountryCode,
		&p.NativeLang, &p.TargetLang, &p.Level, &p.About, &p.Timezone, &p.CreatedAt,
	)
	if err != nil {
		r.log.Error("GetProfile: query failed", logger.Error(err), logger.String("user_id", userID))
		return nil, err
	}
	return &p, nil
}

func (r *profileRepo) UpdateProfile(ctx context.Context, userID string, req models.UpdateProfileRequest) error {
	sets := make([]string, 0, 12)
	args := make([]any, 0, 12)
	i := 1
	add := func(col string, v any) {
		sets = append(sets, fmt.Sprintf("%s=$%d", col, i))
		args = append(args, v)
		i++
	}

	if req.DisplayName != nil {
		add("display_name", *req.DisplayName)
	}
	if req.AvatarURL != nil {
		add("avatar_url", *req.AvatarURL)
	}
	if req.Age != nil {
		add("age", *req.Age)
	}
	if req.Gender != nil {
		add("gender", *req.Gender)
	}
	if req.CountryCode != nil {
		cc := strings.ToUpper(strings.TrimSpace(*req.CountryCode))
		add("country_code", cc)
	}
	if req.NativeLang != nil {
		add("native_lang", *req.NativeLang)
	}
	if req.TargetLang != nil {
		add("target_lang", *req.TargetLang)
	}
	if req.Level != nil {
		add("level", *req.Level)
	}
	if req.About != nil {
		add("about", *req.About)
	}
	if req.Timezone != nil {
		add("timezone", *req.Timezone)
	}

	if len(sets) == 0 {
		return nil // nothing to update
	}
	q := fmt.Sprintf(`UPDATE users SET %s WHERE id=$%d`, strings.Join(sets, ", "), i)
	args = append(args, userID)

	_, err := r.db.Exec(ctx, q, args...)
	if err != nil {
		r.log.Error("UpdateProfile: exec failed", logger.Error(err), logger.String("user_id", userID))
		return err
	}
	return nil
}
