package usecase

import (
	"bytes"
	"errors"
	"github.com/rodrigoachilles/simple-weather-otel/weather-service/internal/dto"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	mockAPIKey              = "your_api_key"
	mockLocale              = "Sao%20Paulo"
	mockInvalidLocale       = "Invalid%20Locale"
	mockWeatherResponseBody = `{"location": {"name": "Sao Paulo"}, "current": {}}`
)

func TestWeatherFinder_Execute_Success(t *testing.T) {
	mockRoundTripper := new(MockRoundTripper)
	mockClient := &http.Client{Transport: mockRoundTripper}

	mockRoundTripper.On("RoundTrip", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(mockWeatherResponseBody))),
	}, nil)

	origAPIKey := os.Getenv(keyWeatherApi)
	_ = os.Setenv(keyWeatherApi, mockAPIKey)
	defer func(key, value string) {
		_ = os.Setenv(key, value)
	}(keyWeatherApi, origAPIKey)

	finder := NewWeatherFinder(mockClient)
	output, err := finder.Execute(mockLocale)

	assert.Nil(t, err)
	assert.NotNil(t, output)
}

func TestWeatherFinder_Execute_HttpClientError(t *testing.T) {
	mockRoundTripper := new(MockRoundTripper)
	mockClient := &http.Client{Transport: mockRoundTripper}

	mockRoundTripper.On("RoundTrip", mock.Anything).Return(nil, errors.New("mocked http client error"))

	origAPIKey := os.Getenv(keyWeatherApi)
	_ = os.Setenv(keyWeatherApi, mockAPIKey)
	defer func(key, value string) {
		_ = os.Setenv(key, value)
	}(keyWeatherApi, origAPIKey)

	finder := NewWeatherFinder(mockClient)
	_, err := finder.Execute(mockLocale)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "mocked http client error")
}

func TestWeatherFinder_Execute_InvalidApiKey(t *testing.T) {
	mockRoundTripper := new(MockRoundTripper)
	mockClient := &http.Client{Transport: mockRoundTripper}

	mockRoundTripper.On("RoundTrip", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusUnauthorized,
		Body:       io.NopCloser(bytes.NewReader([]byte(""))),
	}, nil)

	origAPIKey := os.Getenv(keyWeatherApi)
	_ = os.Setenv(keyWeatherApi, mockAPIKey)
	defer func(key, value string) {
		_ = os.Setenv(key, value)
	}(keyWeatherApi, origAPIKey)

	finder := NewWeatherFinder(mockClient)
	_, err := finder.Execute(mockLocale)

	assert.NotNil(t, err)
	assert.Equal(t, "API key is invalid or not provided", err.Error())
}

func TestWeatherFinder_Execute_InvalidResponse(t *testing.T) {
	mockRoundTripper := new(MockRoundTripper)
	mockClient := &http.Client{Transport: mockRoundTripper}

	expectedOutput := &dto.WeatherOutput{}

	mockRoundTripper.On("RoundTrip", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"error": {"code": 1006, "message": "No matching location found."}}`))),
	}, nil)

	origAPIKey := os.Getenv(keyWeatherApi)
	_ = os.Setenv(keyWeatherApi, mockAPIKey)
	defer func(key, value string) {
		_ = os.Setenv(key, value)
	}(keyWeatherApi, origAPIKey)

	finder := NewWeatherFinder(mockClient)
	output, err := finder.Execute(mockInvalidLocale)

	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, output)
}
