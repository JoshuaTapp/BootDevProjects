package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type apiConfig struct {
	fileserverHits int
}

type chirp struct{
	Body string `json:"body"`
	Id	string `json:"id"`
}

func main() {
	const port = "8080"
	cfg := new(apiConfig)

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("."))

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("GET /api/reset", cfg.resetHandler)
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("POST /api/validate_chirp", isValidChirpHandler)
	mux.HandleFunc("GET /api/chirps" getChirpHandler)
	mux.HandleFunc("POST /api/chirps" postChirpHandler)

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

func isValidChirpHandler(w http.ResponseWriter, r *http.Request) {

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

func postChirpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type chirp struct {
		Body string `json:"body"`
	}

	// decode chirp
	decoder := json.NewDecoder(r.Body)
	c := chirp{}
	err := decoder.Decode(&c)
	if err != nil {
		log.Println("Error decoding params: ", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(c.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	validRespBody := chirp{
		Body: removeProfanity(c.Body),
		Id: 1,
	}

	respondWithJSON(w, http.StatusCreated, validRespBody)
}

func getChirpHandler(w http.ResponseWriter, r *http.Request) {
	var chirps []chirp

	
}