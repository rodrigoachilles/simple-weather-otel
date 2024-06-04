package usecase

import (
	"encoding/json"
	"fmt"
	"github.com/rodrigoachilles/simple-weather-otel/cep-service/internal/dto"
	"github.com/rodrigoachilles/simple-weather/pkg/log"
	"io"
	"net/http"
)

const urlViacepApi = "https://viacep.com.br/ws/%s/json/"

type LocaleFinder struct {
	httpClient *http.Client
}

func NewLocaleFinder(httpClient *http.Client) *LocaleFinder {
	return &LocaleFinder{httpClient: httpClient}
}

func (l *LocaleFinder) Execute(cep string) (*dto.LocaleOutput, error) {
	requestURL := fmt.Sprintf(urlViacepApi, cep)

	log.Logger.Debug().Msg(fmt.Sprintf("Calling api url: %s", requestURL))

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	res, err := l.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	_ = res.Body.Close()

	log.Logger.Debug().Msg(fmt.Sprintf("Response body: %s", string(body)))

	var output dto.LocaleOutput
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
