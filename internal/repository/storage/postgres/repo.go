package postgres

import (
	"context"
	"net/url"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/exp/slog"

	entity "github.com/sreway/shorturl/internal/domain/url"
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

func (r *repo) Add(ctx context.Context, id, userID [16]byte, value *url.URL) error {
	// need implement this
	return nil
}

func (r *repo) Get(ctx context.Context, id [16]byte) (value url.URL, userID [16]byte, err error) {
	// need implement this
	return url.URL{}, [16]byte{}, nil
}

func (r *repo) GetByUserID(ctx context.Context, userID [16]byte) ([]entity.URL, error) {
	// need implement this
	return nil, nil
}

func (r *repo) Close() error {
	r.pool.Close()
	return nil
}

func New(ctx context.Context, dsn string) (*repo, error) {
	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("repository", "postgres")}))

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &repo{
		pool:   pool,
		logger: log,
	}, nil
}
