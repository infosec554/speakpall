// storage/postgres/profile_repo.go
package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
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

// -------------------- PROFILE --------------------

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

// -------------------- INTERESTS --------------------

func (r *profileRepo) GetUserInterests(ctx context.Context, userID string) ([]int, error) {
	const q = `SELECT interest_id FROM user_interests WHERE user_id=$1 ORDER BY interest_id`
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		r.log.Error("GetUserInterests: query failed", logger.Error(err), logger.String("user_id", userID))
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *profileRepo) ReplaceUserInterests(ctx context.Context, userID string, interestIDs []int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, `DELETE FROM user_interests WHERE user_id=$1`, userID); err != nil {
		r.log.Error("ReplaceUserInterests: delete failed", logger.Error(err), logger.String("user_id", userID))
		return err
	}
	for _, id := range interestIDs {
		if _, err := tx.Exec(ctx,
			`INSERT INTO user_interests(user_id, interest_id) VALUES ($1,$2)`,
			userID, id,
		); err != nil {
			r.log.Error("ReplaceUserInterests: insert failed", logger.Error(err), logger.String("user_id", userID))
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

// -------------------- SETTINGS --------------------

func (r *profileRepo) GetUserSettings(ctx context.Context, userID string) (*models.UserSettings, error) {
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
			// jadvaldagi defaultlar
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

func (r *profileRepo) UpsertUserSettings(ctx context.Context, userID string, req models.UpdateSettingsRequest) error {
	// hozirgi qiymatlarni o‘qiymiz (yo‘q bo‘lsa default)
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

// -------------------- MATCH PREFERENCES --------------------

func (r *profileRepo) GetMatchPrefs(ctx context.Context, userID string) (*models.MatchPreferences, error) {
	const q = `
SELECT target_lang, min_level, max_level, gender_filter, min_rating, countries_allow,
       COALESCE(to_char(updated_at,'YYYY-MM-DD"T"HH24:MI:SS"Z"'),'') AS updated_at
FROM match_preferences WHERE user_id=$1`
	var mp models.MatchPreferences
	var countries []string
	err := r.db.QueryRow(ctx, q, userID).Scan(
		&mp.TargetLang, &mp.MinLevel, &mp.MaxLevel, &mp.GenderFilter, &mp.MinRating, &countries, &mp.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			// bo'sh preferensiya qaytaramiz
			return &models.MatchPreferences{}, nil
		}
		r.log.Error("GetMatchPrefs: query failed", logger.Error(err), logger.String("user_id", userID))
		return nil, err
	}
	if len(countries) > 0 {
		mp.CountriesAllow = countries
	}
	return &mp, nil
}

func (r *profileRepo) UpsertMatchPrefs(ctx context.Context, userID string, req models.UpdateMatchPrefsRequest) error {
	// min/max validatsiya
	if req.MinLevel != nil && req.MaxLevel != nil && *req.MinLevel > *req.MaxLevel {
		return fmt.Errorf("min_level cannot be greater than max_level")
	}

	// hozirgi qiymatlarni o‘qiymiz
	cur, _ := r.GetMatchPrefs(ctx, userID)

	if req.TargetLang != nil {
		cur.TargetLang = req.TargetLang
	}
	if req.MinLevel != nil {
		cur.MinLevel = req.MinLevel
	}
	if req.MaxLevel != nil {
		cur.MaxLevel = req.MaxLevel
	}
	if req.GenderFilter != nil {
		cur.GenderFilter = req.GenderFilter
	}
	if req.MinRating != nil {
		cur.MinRating = req.MinRating
	}
	if req.CountriesAllow != nil {
		cur.CountriesAllow = req.CountriesAllow
	}

	const q = `
INSERT INTO match_preferences (user_id, target_lang, min_level, max_level, gender_filter, min_rating, countries_allow, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7, now())
ON CONFLICT (user_id) DO UPDATE SET
  target_lang     = EXCLUDED.target_lang,
  min_level       = EXCLUDED.min_level,
  max_level       = EXCLUDED.max_level,
  gender_filter   = EXCLUDED.gender_filter,
  min_rating      = EXCLUDED.min_rating,
  countries_allow = EXCLUDED.countries_allow,
  updated_at      = now()`
	_, err := r.db.Exec(
		ctx, q, userID,
		cur.TargetLang, cur.MinLevel, cur.MaxLevel, cur.GenderFilter, cur.MinRating, cur.CountriesAllow,
	)
	if err != nil {
		r.log.Error("UpsertMatchPrefs: exec failed", logger.Error(err), logger.String("user_id", userID))
		return err
	}
	return nil
}
