package postgres

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	"speakpall/config"
	"speakpall/pkg/logger"
	"speakpall/storage"
)

type Store struct {
	pool  *pgxpool.Pool
	log   logger.ILogger
	redis storage.IRedisStorage
}

func New(ctx context.Context, cfg config.Config, log logger.ILogger, redis storage.IRedisStorage) (storage.IStorage, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Error("error while parsing config", logger.Error(err))
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Error("error while connecting to database", logger.Error(err))
		return nil, err
	}

	absPath, err := filepath.Abs("migrations/postgres")
	if err != nil {
		log.Error("failed to get absolute path for migrations", logger.Error(err))
		return nil, err
	}
	m, err := migrate.New("file://"+absPath, url)
	if err != nil {
		log.Error("migration error", logger.Error(err))
		return nil, err
	}
	if err = m.Up(); err != nil && !strings.Contains(err.Error(), "no change") {
		log.Error("migration up error", logger.Error(err))
		return nil, err
	}

	log.Info("postgres connected and migrated")

	return &Store{
		pool:  pool,
		log:   log,
		redis: redis,
	}, nil
}

func (s *Store) Close() {
	s.pool.Close()
}

func (s *Store) User() storage.IUserStorage {
	return NewUserRepo(s.pool, s.log)
}

func (s *Store) Profile()  storage.IProfileStorage {
	return NewProfileRepo(s.pool, s.log)
}

func (s *Store) Redis() storage.IRedisStorage {
	return s.redis
}
