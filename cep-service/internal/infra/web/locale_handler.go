package web

import (
	"encoding/json"
	"fmt"
	"github.com/rodrigoachilles/simple-weather-otel/cep-service/internal/dto"
	"github.com/rodrigoachilles/simple-weather-otel/cep-service/internal/usecase"
	"github.com/rodrigoachilles/simple-weather-otel/pkg/log"
	"net/http"
)

type LocaleHandler struct {
	localeFinder usecase.Finder
}

func NewLocaleHandler(localeFinder usecase.Finder) *LocaleHandler {
	return &LocaleHandler{
		localeFinder: localeFinder,
	}
}

func (h *LocaleHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input dto.LocaleInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := struct {
			Message string `json:"message"`
		}{
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(msg)
		return
	}

	if len(input.Cep) != 8 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(&dto.ErrorOutput{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "invalid zipcode",
		})
		return
	}

	output, err := h.localeFinder.Execute(input.Cep)
	if err != nil {
		log.Logger.Error().Err(err).Msg(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(&dto.ErrorOutput{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
		return
	}

	if output.Localidade == "" {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(&dto.ErrorOutput{
			StatusCode: http.StatusNotFound,
			Message:    "can not find zipcode",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Logger.Info().Msg(fmt.Sprintf("%s", output))
	_ = json.NewEncoder(w).Encode(output)
}
