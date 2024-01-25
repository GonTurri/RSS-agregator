package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GonTurri/RSS-agregator/internal/database"
	"github.com/google/uuid"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handlerErr(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (apiCfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Name string
	}
	decoder := json.NewDecoder(r.Body)
	param := params{}
	err := decoder.Decode(&param)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error parsing json: %v", err))
		return
	}
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      param.Name,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating user: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, dbUserToUser(user))
}
