package inmemory

import (
	"context"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"shortener/internal/entities"
	"sync"
)

type InMemoryStorage struct {
	log  *zap.SugaredLogger
	data map[string]string
	mu   *sync.Mutex
}

func NewInMemoryStorage(log *zap.SugaredLogger) (*InMemoryStorage, error) {
	if log == nil {
		return nil, errors.WithMessage(entities.ErrInvalidParam, "empty logger")
	}

	return &InMemoryStorage{
		log:  log,
		data: make(map[string]string),
		mu:   &sync.Mutex{},
	}, nil
}

func (in *InMemoryStorage) Close() {
}

func (in *InMemoryStorage) CreateShortLink(_ context.Context, shortLink string, originalLink string) error {
	in.mu.Lock()
	defer in.mu.Unlock()

	if _, ok := in.data[shortLink]; ok {
		in.log.Error(entities.ErrAlreadyExists)

		//return errors.WithMessage(entities.ErrInternal, "link already exist")
		return entities.ErrAlreadyExists
	}

	in.data[shortLink] = originalLink

	return nil
}
func (in *InMemoryStorage) GetOriginalLink(_ context.Context, shortLink string) (string, error) {
	originalLink, ok := in.data[shortLink]
	if !ok {
		in.log.Error(entities.ErrNotFound)

		return "", errors.WithMessage(entities.ErrNotFound, "no original link for this query")
	}

	return originalLink, nil
}
