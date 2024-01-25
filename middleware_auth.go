package main

import (
	"net/http"

	"github.com/GonTurri/RSS-agregator/internal/auth"
	"github.com/GonTurri/RSS-agregator/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)

		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)

		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		handler(w, r, user)
	}
}
