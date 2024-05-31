package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/JoshuaTapp/BootDevProjects/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
}

var (
	dbConn *database.DB
)

func main() {
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	dbConn, _ = database.NewDB("database.json", *dbg)

	const port = "8080"
	cfg := new(apiConfig)

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("."))

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("GET /api/reset", cfg.resetHandler)
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /api/chirps", getChirpHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", getChirpHandler)
	mux.HandleFunc("POST /api/chirps", postChirpHandler)

	mux.HandleFunc("POST /api/users", postUsersHandler)

	loggingHandler := loggingMiddleware(mux)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: loggingHandler,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Metrics Called - Hits: ", cfg.fileserverHits)

	tmpl := fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>
	`, cfg.fileserverHits)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tmpl))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Resetting hits...")
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("middleware triggered - Prev Hits: ", cfg.fileserverHits)
		cfg.fileserverHits++
		log.Print("New Hits: ", cfg.fileserverHits)
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Visited page: %s\n", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func validateChirp(msg string) (bool, error) {
	if len(msg) > 140 {
		return false, errors.New("Chirp is too long")
	}

	return true, nil
}

func removeProfanity(msg string) string {
	theProfane := map[string]string{
		"kerfuffle": "****",
		"sharbert":  "****",
		"fornax":    "****",
	}

	strSlice := strings.Split(msg, " ")
	log.Print("strSlice: ", strSlice)

	for i, s := range strSlice {
		lcStr := strings.ToLower(s)
		if v, ok := theProfane[lcStr]; ok {
			strSlice[i] = v
		}
	}

	return strings.Join(strSlice, " ")

}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	respBody := struct {
		Error string `json:"error"`
	}{
		Error: msg,
	}

	rb, err := json.Marshal(&respBody)
	if err != nil {
		log.Printf("Error mashalling JSON: %s", err)
		return
	}
	w.Write(rb)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	rb, err := json.Marshal(&payload)
	if err != nil {
		log.Printf("Error mashalling JSON: %s", err)
		return
	}
	w.Write(rb)
}

func decodeJSON(r *http.Request, shape interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(shape)
}

func postChirpHandler(w http.ResponseWriter, r *http.Request) {
	c := &struct {
		Body string `json:"body"`
	}{}
	err := decodeJSON(r, c)
	if err != nil {
		log.Println("Error decoding params: ", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if ok, err := validateChirp(c.Body); !ok {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := dbConn.CreateChirp(removeProfanity(c.Body))
	if err != nil {
		log.Print("failed to create chirp", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, chirp)
}

func getChirpHandler(w http.ResponseWriter, r *http.Request) {
	var chirps []database.Chirp

	chirps, err := dbConn.GetChirps()
	if err != nil {
		log.Print("failure getting chirps", err)
	}

	id, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		respondWithJSON(w, http.StatusOK, chirps)
		return
	}

	c, err := dbConn.GetChirp(id)
	if err != nil {
		respondWithError(w, 404, "chirp not found")
		return
	}
	respondWithJSON(w, http.StatusOK, c)
}

func postUsersHandler(w http.ResponseWriter, r *http.Request) {
	u := &struct {
		Email string `json:"email"`
	}{}

	err := decodeJSON(r, u)
	if err != nil {
		log.Println("Error decoding params: ", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	user, err := dbConn.CreateUser(u.Email)
	if err != nil {
		log.Print("failed to create user", err)
		return
	}
	respondWithJSON(w, 201, user)
}
