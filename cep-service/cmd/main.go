package main

import (
	"fmt"
	"github.com/rodrigoachilles/simple-weather-otel/cep-service/internal/infra/web"
	"github.com/rodrigoachilles/simple-weather-otel/cep-service/internal/usecase"
	"github.com/rodrigoachilles/simple-weather-otel/pkg/log"
	"net/http"
)

func main() {
	port := ":8080"
	log.Logger.Info().Msg(fmt.Sprintf("Starting server on port %s ...", port[1:]))

	mux := http.NewServeMux()
	localeFinder := usecase.NewLocaleFinder(http.DefaultClient)
	mux.HandleFunc("POST /", web.NewLocaleHandler(localeFinder).Handle)

	log.Logger.Fatal().Err(http.ListenAndServe(port, mux))
}
