package main

import (
	"context"
	"errors"
	"fmt"
	otel "github.com/rodrigoachilles/simple-weather-otel"
	"github.com/rodrigoachilles/simple-weather-otel/cep-service/internal/infra/web"
	"github.com/rodrigoachilles/simple-weather-otel/cep-service/internal/usecase"
	"github.com/rodrigoachilles/simple-weather-otel/pkg/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Logger.Fatal().Err(err).Msg(err.Error())
	}
}

func run() (err error) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	otelShutdown, err := otel.SetupOTelSDK("cep-service", ctx)
	if err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	serverPort := os.Getenv("CEP_SERVICE_SERVER_PORT")
	if serverPort == "" {
		serverPort = ":8080"
	}
	srv := &http.Server{
		Addr:         serverPort,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      newHTTPHandler(),
	}
	srvErr := make(chan error, 1)
	go func() {
		log.Logger.Info().Msg(fmt.Sprintf("Starting server on port '%s'...", serverPort[1:]))
		srvErr <- srv.ListenAndServe()
	}()

	select {
	case err = <-srvErr:
		return
	case <-ctx.Done():
		stop()
	}

	err = srv.Shutdown(context.Background())
	return
}

func newHTTPHandler() http.Handler {
	mux := http.NewServeMux()

	localeFinder := usecase.NewLocaleFinder(http.DefaultClient)
	mux.HandleFunc("POST /", web.NewLocaleHandler(localeFinder).Handle)

	return otelhttp.NewHandler(mux, "/")
}
