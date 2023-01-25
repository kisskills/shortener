package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"shortener/internal/entities"
)

type PGStorage struct {
	log    *zap.SugaredLogger
	cancel context.CancelFunc
	db     *pgxpool.Pool
}

func NewPGStorage(log *zap.SugaredLogger, dsn string) (*PGStorage, error) {
	if log == nil {
		return nil, errors.WithMessage(entities.ErrInvalidParam, "empty logger")
	}

	if dsn == "" {
		return nil, errors.WithMessage(entities.ErrInvalidParam, "empty dsn")
	}

	st := &PGStorage{
		log: log,
	}

	ctx, cancel := context.WithCancel(context.Background())
	st.cancel = cancel

	conn, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	st.db = conn

	return st, nil
}

func (s *PGStorage) Close() {
	s.db.Close()
	s.cancel()
}

func (s *PGStorage) CreateShortLink(ctx context.Context, shortLink string, originalLink string) error {
	query := `INSERT INTO ozon.links (short_link, original_link) VALUES ($1, $2)`

	_, err := s.db.Exec(ctx, query, shortLink, originalLink)
	if err != nil {
		s.log.Error(err)

		duplicateEntryError := &pgconn.PgError{Code: "23505"}
		if errors.As(err, &duplicateEntryError) {
			return entities.ErrAlreadyExists
		}
		return errors.WithMessage(entities.ErrInternal, err.Error())
	}

	return nil
}

func (s *PGStorage) GetOriginalLink(ctx context.Context, shortLink string) (string, error) {
	query := `SELECT original_link FROM ozon.links WHERE short_link = $1`

	var originalLink string

	row := s.db.QueryRow(ctx, query, &shortLink)
	err := row.Scan(&originalLink)
	if errors.Is(err, pgx.ErrNoRows) {
		err = errors.WithMessage(entities.ErrNotFound, "no original link for this query")
		s.log.Error(err)
		return "", err
	}

	if err != nil {
		s.log.Error(err)
		return "", errors.WithMessage(entities.ErrInternal, err.Error())
	}

	return originalLink, nil
}
