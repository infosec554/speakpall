package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"speakpall/pkg/logger"
	"speakpall/storage"
)

type interesRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewInteresRepo(db *pgxpool.Pool, log logger.ILogger) storage.IUserInterestsStorage {
	return &interesRepo{db: db, log: log}
}

func (r *interesRepo) GetUserInterests(ctx context.Context, userID string) ([]int, error) {
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

func (r *interesRepo) ReplaceUserInterests(ctx context.Context, userID string, interestIDs []int) error {
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
