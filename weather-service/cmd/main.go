package main

import (
	"fmt"
	"github.com/rodrigoachilles/simple-weather-otel/pkg/log"
	"github.com/rodrigoachilles/simple-weather-otel/weather-service/internal/infra/web"
	"github.com/rodrigoachilles/simple-weather-otel/weather-service/internal/usecase"
	"net/http"
)

func main() {
	port := ":8081"
	log.Logger.Info().Msg(fmt.Sprintf("Starting server on port %s ...", port[1:]))

	mux := http.NewServeMux()
	weatherFinder := usecase.NewWeatherFinder(http.DefaultClient)
	localeFinder := usecase.NewLocaleFinder(http.DefaultClient)
	mux.HandleFunc("GET /{cep}", web.NewWeatherHandler(weatherFinder, localeFinder).Handle)

	log.Logger.Fatal().Err(http.ListenAndServe(port, mux))
}
