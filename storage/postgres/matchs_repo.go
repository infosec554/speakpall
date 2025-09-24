package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"speakpall/api/models"
	"speakpall/pkg/logger"
	"speakpall/storage"
)

type matchsRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewMatchsRepo(db *pgxpool.Pool, log logger.ILogger) storage.IMatchPreferencesStorage {
	return &matchsRepo{db: db, log: log}
}
func (r *matchsRepo) GetMatchPrefs(ctx context.Context, userID string) (*models.MatchPreferences, error) {
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

func (r *matchsRepo) UpsertMatchPrefs(ctx context.Context, userID string, req models.UpdateMatchPrefsRequest) error {
	if req.MinLevel != nil && req.MaxLevel != nil && *req.MinLevel > *req.MaxLevel {
		return fmt.Errorf("min_level cannot be greater than max_level")
	}

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
