package postgres

import (
	"context"
	"errors"
	"net/url"
	"os"

	"github.com/golang-migrate/migrate/v4"
	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/exp/slog"

	"github.com/sreway/shorturl/internal/config"
	entity "github.com/sreway/shorturl/internal/domain/url"
	"github.com/sreway/shorturl/internal/usecases/shortener"
)

type repo struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func (r *repo) Ping(ctx context.Context) error {
	if err := r.pool.Ping(ctx); err != nil {
		r.logger.Error("failed execute empty sql statement", err, slog.String("func", "Ping"))
		return err
	}
	return nil
}

func (r *repo) Add(ctx context.Context, item entity.URL) error {
	query := "INSERT INTO urls (id, user_id, original_url) VALUES ($1, $2, $3)"
	_, err := r.pool.Exec(ctx, query, item.ID(), item.UserID(), item.LongURL().String())
	if err != nil {
		r.logger.Error("failed insert short url", err)
		return err
	}
	return nil
}

func (r *repo) Get(ctx context.Context, id [16]byte) (entity.URL, error) {
	var (
		userID [16]byte
		rawURL string
	)
	query := "SELECT user_id, original_url FROM urls WHERE id = $1"
	err := r.pool.QueryRow(ctx, query, id).Scan(&userID, &rawURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, shortener.ErrNotFound
		}
		return nil, err
	}

	value, err := url.ParseRequestURI(rawURL)
	if err != nil {
		r.logger.Error("failed parse raw url", err, slog.String("func", "Get"),
			slog.String("url", rawURL))
		return nil, err
	}

	return entity.NewURL(id, userID, nil, value), nil
}

func (r *repo) GetByUserID(ctx context.Context, userID [16]byte) ([]entity.URL, error) {
	urls := make([]entity.URL, 0)
	query := "SELECT id, original_url FROM urls WHERE user_id = $1"
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id     [16]byte
			rawURL string
		)
		if err = rows.Scan(&id, &rawURL); err != nil {
			return nil, err
		}

		value, err := url.ParseRequestURI(rawURL)
		if err != nil {
			r.logger.Error("failed parse raw url", err, slog.String("func", "Get"),
				slog.String("url", rawURL))
			return nil, err
		}

		urls = append(urls, entity.NewURL(id, userID, nil, value))
	}

	return urls, nil
}

func (r *repo) Close() error {
	r.pool.Close()
	return nil
}

func (r *repo) Batch(ctx context.Context, urls []entity.URL) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	if err != nil {
		return err
	}

	query := "INSERT INTO urls (id, user_id, original_url) VALUES ($1, $2, $3)"

	for _, item := range urls {
		_, err = r.pool.Exec(ctx, query, item.ID(), item.UserID(), item.LongURL().String())
		if err != nil {
			r.logger.Error("failed batch insert short url", err)
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *repo) migrate(migrateURL string) error {
	m, err := migrate.New(migrateURL, r.pool.Config().ConnConfig.ConnString())
	defer func() {
		_, _ = m.Close()
	}()

	if err != nil {
		return err
	}
	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		r.logger.Info("no change", slog.String("func", "migrate"))
		return nil
	}

	return err
}

func New(ctx context.Context, config config.Postgres) (*repo, error) {
	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("repository", "postgres")}))

	poolConfig, err := pgxpool.ParseConfig(config.GetDSN())
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	r := &repo{
		pool:   pool,
		logger: log,
	}

	if len(config.GetMigrateURL()) == 0 {
		return r, nil
	}

	err = r.migrate(config.GetMigrateURL())
	if err != nil {
		log.Error("failed apply migrations", err, slog.String("func", "migrate"))
	}

	return r, nil
}
