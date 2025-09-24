package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"speakpall/pkg/logger"
	"speakpall/storage"
)

type friendRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewFriendRepo(db *pgxpool.Pool, log logger.ILogger) storage.IFriendStorage {
	return &friendRepo{db: db, log: log}
}

func (r *friendRepo) AddFriend(ctx context.Context, userID, friendID string) error {
	if userID == friendID {
		return fmt.Errorf("cannot add yourself as friend")
	}
	const q = `
INSERT INTO friends (user_id, friend_user_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(ctx, q, userID, friendID)
	if err != nil {
		r.log.Error("AddFriend: exec failed", logger.Error(err), logger.String("user_id", userID), logger.String("friend_id", friendID))
		return err
	}
	return nil
}

func (r *friendRepo) RemoveFriend(ctx context.Context, userID, friendID string) error {
	const q = `DELETE FROM friends WHERE user_id=$1 AND friend_user_id=$2`
	_, err := r.db.Exec(ctx, q, userID, friendID)
	if err != nil {
		r.log.Error("RemoveFriend: exec failed", logger.Error(err), logger.String("user_id", userID), logger.String("friend_id", friendID))
		return err
	}
	return nil
}

func (r *friendRepo) ListFriends(ctx context.Context, userID string) ([]string, error) {
	const q = `SELECT friend_user_id FROM friends WHERE user_id=$1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		r.log.Error("ListFriends: query failed", logger.Error(err), logger.String("user_id", userID))
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return ids, nil
}

func (r *friendRepo) IsFriend(ctx context.Context, userID, friendID string) (bool, error) {
	const q = `SELECT 1 FROM friends WHERE user_id=$1 AND friend_user_id=$2 LIMIT 1`
	var tmp int
	err := r.db.QueryRow(ctx, q, userID, friendID).Scan(&tmp)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		r.log.Error("IsFriend: query failed", logger.Error(err), logger.String("user_id", userID), logger.String("friend_id", friendID))
		return false, err
	}
	return true, nil
}
