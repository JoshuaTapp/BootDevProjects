package main

import (
	"encoding/json"
	"net/http"

	"github.com/JoshuaTapp/BootDevProjects/chirpy/internal/auth"
	"github.com/JoshuaTapp/BootDevProjects/chirpy/internal/database"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type polkaRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetPolkaKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no api key provided")
		return
	}
	if apiKey != cfg.polkaSecret {
		respondWithError(w, http.StatusUnauthorized, "invalid polka key")
		return
	}

	decoder := json.NewDecoder(r.Body)
	payload := polkaRequest{}
	err = decoder.Decode(&payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if payload.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.DB.UpgradeUser(payload.Data.UserID)
	if err != nil {
		if err == database.ErrNotExist {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)

}
