package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:       dbChirp.ID,
		Body:     dbChirp.Body,
		AuthorID: dbChirp.AuthorID,
	})
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("author_id"))

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		if err == nil && dbChirp.AuthorID != id {
			continue
		}

		chirps = append(chirps, Chirp{
			ID:       dbChirp.ID,
			Body:     dbChirp.Body,
			AuthorID: dbChirp.AuthorID,
		})
	}

	sortType := r.URL.Query().Get("sort")
	sort.Slice(chirps, func(i, j int) bool {
		if sortType == "desc" {
			return chirps[i].ID >= chirps[j].ID
		} else {
			return chirps[i].ID < chirps[j].ID
		}
	})

	respondWithJSON(w, http.StatusOK, chirps)
}
