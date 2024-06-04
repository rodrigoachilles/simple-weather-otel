package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rodrigoachilles/simple-weather-otel/weather-service/internal/dto"
	"github.com/rodrigoachilles/simple-weather/pkg/log"
	"io"
	"net/http"
)

const urlLocaleApi = "http://localhost:8080/"

type LocaleFinder struct {
	httpClient *http.Client
}

func NewLocaleFinder(httpClient *http.Client) *LocaleFinder {
	return &LocaleFinder{httpClient: httpClient}
}

func (l *LocaleFinder) Execute(cep string) (interface{}, error) {
	input := &dto.LocaleInput{
		Cep: cep,
	}
	inputJson, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	log.Logger.Debug().Msg(fmt.Sprintf("Calling api url: %s", urlLocaleApi))

	req, err := http.NewRequest(http.MethodPost, urlLocaleApi, bytes.NewBuffer(inputJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	res, err := l.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	_ = res.Body.Close()

	var output dto.LocaleOutput
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
