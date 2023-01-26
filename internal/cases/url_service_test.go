package cases

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shortener/internal/cases/mock"
	"shortener/internal/entities"
	"shortener/internal/proto_gen/shortener"
	"testing"
)

func TestUrlService_CreateShortLink(t *testing.T) {
	type mockActions func(r *mock.MockStorage, shortLink, originalLink string)

	tests := []struct {
		name        string
		mockActions mockActions
		args        *shortener.NewLinkRequest
		wantResp    *shortener.NewLinkResponse
		expectedErr error
	}{
		{
			name: "OK_1",
			mockActions: func(r *mock.MockStorage, shortLink, originalLink string) {
				r.EXPECT().CreateShortLink(context.Background(), shortLink, originalLink).Return(nil)
			},
			args:        &shortener.NewLinkRequest{Link: "google.com"},
			wantResp:    &shortener.NewLinkResponse{ShortLink: "kwUIdA9TOq"},
			expectedErr: nil,
		},
		{
			name: "OK_2",
			mockActions: func(r *mock.MockStorage, shortLink, originalLink string) {
				r.EXPECT().CreateShortLink(context.Background(), shortLink, originalLink).Return(nil)
			},
			args:        &shortener.NewLinkRequest{Link: "ozon.ru"},
			wantResp:    &shortener.NewLinkResponse{ShortLink: "naNLoN5x10"},
			expectedErr: nil,
		},
		{
			name:        "Error empty link",
			mockActions: func(r *mock.MockStorage, shortLink, originalLink string) {},
			args:        &shortener.NewLinkRequest{Link: ""},
			wantResp:    nil,
			expectedErr: status.Error(codes.InvalidArgument,
				errors.WithMessage(entities.ErrInvalidParam, "empty original link").Error()),
		},
		{
			name: "Error already exist 1",
			mockActions: func(r *mock.MockStorage, shortLink, originalLink string) {
				r.EXPECT().CreateShortLink(context.Background(), shortLink, originalLink).Return(entities.ErrAlreadyExists)
			},
			args:        &shortener.NewLinkRequest{Link: "github.com"},
			wantResp:    nil,
			expectedErr: status.Error(codes.AlreadyExists, entities.ErrAlreadyExists.Error()),
		},
		{
			name: "Error already exist 2",
			mockActions: func(r *mock.MockStorage, shortLink, originalLink string) {
				r.EXPECT().CreateShortLink(context.Background(), shortLink, originalLink).Return(entities.ErrAlreadyExists)
			},
			args:        &shortener.NewLinkRequest{Link: "job.ozon.ru"},
			wantResp:    nil,
			expectedErr: status.Error(codes.AlreadyExists, entities.ErrAlreadyExists.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mock.NewMockStorage(ctrl)
			tt.mockActions(storage, Hasher(tt.args.Link), tt.args.Link)

			logger, _ := zap.NewProduction()
			svc, _ := NewUrlService(logger.Sugar(), storage)

			resp, err := svc.CreateShortLink(context.Background(), tt.args)

			assert.Equal(t, tt.wantResp, resp)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestUrlService_GetOriginalLink(t *testing.T) {
	type mockActions func(r *mock.MockStorage, shortLink string)

	tests := []struct {
		name        string
		mockActions mockActions
		args        *shortener.GetLinkRequest
		wantResp    *shortener.GetLinkResponse
		expectedErr error
	}{
		{
			name: "OK_1",
			args: &shortener.GetLinkRequest{Link: "kwUIdA9TOq"},
			mockActions: func(r *mock.MockStorage, shortLink string) {
				r.EXPECT().GetOriginalLink(context.Background(), shortLink).Return("google.com", nil)
			},
			wantResp:    &shortener.GetLinkResponse{OrigLink: "google.com"},
			expectedErr: nil,
		},
		{
			name: "OK_2",
			args: &shortener.GetLinkRequest{Link: "naNLoN5x10"},
			mockActions: func(r *mock.MockStorage, shortLink string) {
				r.EXPECT().GetOriginalLink(context.Background(), shortLink).Return("ozon.ru", nil)
			},
			wantResp:    &shortener.GetLinkResponse{OrigLink: "ozon.ru"},
			expectedErr: nil,
		},
		{
			name:        "Error wrong short link(too short)",
			args:        &shortener.GetLinkRequest{Link: "012"},
			mockActions: func(r *mock.MockStorage, shortLink string) {},
			wantResp:    nil,
			expectedErr: status.Error(codes.InvalidArgument,
				errors.WithMessage(entities.ErrInvalidParam, "wrong short link").Error()),
		},
		{
			name:        "Error wrong short link(too long)",
			args:        &shortener.GetLinkRequest{Link: "kwUIdA9TOqkwUIdA9TOqkwUIdA9TOqkwUIdA9TOq"},
			mockActions: func(r *mock.MockStorage, shortLink string) {},
			wantResp:    nil,
			expectedErr: status.Error(codes.InvalidArgument,
				errors.WithMessage(entities.ErrInvalidParam, "wrong short link").Error()),
		},
		{
			name: "Error not found 1",
			args: &shortener.GetLinkRequest{Link: "0123456789"},
			mockActions: func(r *mock.MockStorage, shortLink string) {
				r.EXPECT().GetOriginalLink(context.Background(), shortLink).Return("",
					errors.WithMessage(entities.ErrNotFound, "no original link for this query"))
			},
			wantResp: nil,
			expectedErr: status.Errorf(codes.NotFound,
				errors.WithMessage(entities.ErrNotFound, "no original link for this query").Error()),
		},
		{
			name: "Error not found 2",
			args: &shortener.GetLinkRequest{Link: "abcdefghjk"},
			mockActions: func(r *mock.MockStorage, shortLink string) {
				r.EXPECT().GetOriginalLink(context.Background(), shortLink).Return("",
					errors.WithMessage(entities.ErrNotFound, "no original link for this query"))
			},
			wantResp: nil,
			expectedErr: status.Errorf(codes.NotFound,
				errors.WithMessage(entities.ErrNotFound, "no original link for this query").Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mock.NewMockStorage(ctrl)
			tt.mockActions(storage, tt.args.Link)

			logger, _ := zap.NewProduction()
			svc, _ := NewUrlService(logger.Sugar(), storage)

			resp, err := svc.GetOriginalLink(context.Background(), tt.args)

			assert.Equal(t, tt.wantResp, resp)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestNewUrlService(t *testing.T) {
	logger, _ := zap.NewProduction()
	tests := []struct {
		name        string
		log         *zap.SugaredLogger
		needStorage bool
		want        *UrlService
		wantErr     error
	}{
		{
			name:        "Empty logger",
			log:         nil,
			needStorage: true,
			want:        nil,
			wantErr:     errors.WithMessage(entities.ErrInvalidParam, "empty logger"),
		},
		{
			name:        "Empty storage",
			log:         logger.Sugar(),
			needStorage: false,
			want:        nil,
			wantErr:     errors.WithMessage(entities.ErrInvalidParam, "empty storage"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var storage Storage

			if tt.needStorage {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				storage = mock.NewMockStorage(ctrl)
			}
			got, err := NewUrlService(tt.log, storage)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
