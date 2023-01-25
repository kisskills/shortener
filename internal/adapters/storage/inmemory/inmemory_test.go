package inmemory

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"shortener/internal/entities"
	"testing"
)

func TestInMemoryStorage_CreateShortLink(t *testing.T) {
	type args struct {
		in0          context.Context
		shortLink    string
		originalLink string
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		expectedErr error
	}{
		{
			name: "OK",
			args: args{
				in0:          context.Background(),
				shortLink:    "google.com",
				originalLink: "google.com",
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "Expected already exists error",
			args: args{
				in0:          context.Background(),
				shortLink:    "google.com",
				originalLink: "google.com",
			},
			wantErr:     true,
			expectedErr: entities.ErrAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, _ := zap.NewProduction()
			in, _ := NewInMemoryStorage(logger.Sugar())
			if tt.wantErr == true {
				_ = in.CreateShortLink(tt.args.in0, tt.args.shortLink, tt.args.originalLink)
			}
			err := in.CreateShortLink(tt.args.in0, tt.args.shortLink, tt.args.originalLink)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestInMemoryStorage_GetOriginalLink(t *testing.T) {
	type args struct {
		in0       context.Context
		shortLink string
	}
	tests := []struct {
		name        string
		args        args
		want        string
		wantErr     bool
		expectedErr error
	}{
		{
			name: "OK",
			args: args{
				in0:       context.Background(),
				shortLink: "google.com",
			},
			want:        "google.com",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "Error Not Found",
			args: args{
				in0:       context.Background(),
				shortLink: "google.com",
			},
			want:        "google.com",
			wantErr:     true,
			expectedErr: errors.WithMessage(entities.ErrNotFound, "no original link for this query"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, _ := zap.NewProduction()
			in, _ := NewInMemoryStorage(logger.Sugar())
			if tt.wantErr == false {
				_ = in.CreateShortLink(tt.args.in0, tt.args.shortLink, tt.want)
			}
			_, err := in.GetOriginalLink(tt.args.in0, tt.args.shortLink)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
