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
			name: "OK",
			mockActions: func(r *mock.MockStorage, shortLink, originalLink string) {
				r.EXPECT().CreateShortLink(context.Background(), shortLink, originalLink).Return(nil)
			},
			args:        &shortener.NewLinkRequest{Link: "google.com"},
			wantResp:    &shortener.NewLinkResponse{ShortLink: "kwUIdA9TOq"},
			expectedErr: nil,
		},
		{
			name:        "Empty link",
			mockActions: func(r *mock.MockStorage, originalLink, shortLink string) {},
			args:        &shortener.NewLinkRequest{Link: ""},
			wantResp:    nil,
			expectedErr: status.Error(codes.InvalidArgument,
				errors.WithMessage(entities.ErrInvalidParam, "empty original link").Error()),
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
			//fmt.Println(resp)
			//fmt.Println(err.Error())

			assert.Equal(t, tt.wantResp, resp)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
