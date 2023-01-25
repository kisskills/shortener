package application

import (
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"shortener/internal/cases/mock"
	"testing"
)

func TestApplication_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := mock.NewMockStorage(ctrl)
	storage.EXPECT().Close()

	logger, _ := zap.NewProduction()
	app := Application{
		log: logger.Sugar(),
		cancel: func() {

		},
		storage: storage,
	}

	app.Stop()
}
