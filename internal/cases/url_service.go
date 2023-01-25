package cases

import (
	"context"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shortener/internal/entities"
	"shortener/internal/proto_gen/shortener"
)

type UrlService struct {
	log *zap.SugaredLogger
	shortener.UnimplementedShortenerServer
	storage Storage
}

func NewUrlService(log *zap.SugaredLogger, storage Storage) (*UrlService, error) {
	if log == nil {
		return nil, errors.WithMessage(entities.ErrInvalidParam, "empty logger")
	}

	if storage == nil || storage == Storage(nil) {
		return nil, errors.WithMessage(entities.ErrInvalidParam, "empty storage")
	}

	return &UrlService{
		log:     log,
		storage: storage,
	}, nil
}

func (s *UrlService) CreateShortLink(ctx context.Context, req *shortener.NewLinkRequest) (*shortener.NewLinkResponse, error) {
	originalLink := req.Link
	if originalLink == "" {
		err := errors.WithMessage(entities.ErrInvalidParam, "empty original link")
		s.log.Error(err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	shortLink := Hasher(originalLink)
	err := s.storage.CreateShortLink(ctx, shortLink, originalLink)
	if err != nil {
		s.log.Error(err)

		if errors.Is(err, entities.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &shortener.NewLinkResponse{ShortLink: shortLink}, nil
}

func (s *UrlService) GetOriginalLink(ctx context.Context, req *shortener.GetLinkRequest) (*shortener.GetLinkResponse, error) {
	shortLink := req.Link
	if len(shortLink) != 10 {
		err := errors.WithMessage(entities.ErrInvalidParam, "wrong short link")
		s.log.Error(err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	originalLink, err := s.storage.GetOriginalLink(ctx, shortLink)
	if errors.Is(err, entities.ErrNotFound) {
		s.log.Error(err)

		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		s.log.Error(err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &shortener.GetLinkResponse{OrigLink: originalLink}, nil
}
