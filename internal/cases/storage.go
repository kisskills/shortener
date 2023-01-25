package cases

import "context"

//go:generate go run github.com/golang/mock/mockgen -destination=./mock/storage.go -source=storage.go -package=mock Storage
type Storage interface {
	CreateShortLink(ctx context.Context, shortLink string, originalLink string) error
	GetOriginalLink(ctx context.Context, shortLink string) (string, error)
	Close()
}
