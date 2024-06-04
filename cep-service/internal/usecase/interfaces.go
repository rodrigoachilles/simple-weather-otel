package usecase

import "github.com/rodrigoachilles/simple-weather-otel/cep-service/internal/dto"

type Finder interface {
	Execute(cep string) (*dto.LocaleOutput, error)
}
