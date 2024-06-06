package usecase

import (
	"context"
	"github.com/rodrigoachilles/simple-weather-otel/cep-service/internal/dto"
)

type Finder interface {
	Execute(ctx context.Context, cep string) (*dto.LocaleOutput, error)
}
