package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/JoshuaTapp/BootDevProjects/chirpy/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type apiConfig struct {
	fileserverHits int
	jwtSecret      []byte
	db             *database.DB
}

var (
	cfg apiConfig
)

const accessTokenInterval int = 60 * 60            // 1 hour
const refreshTokenInterval int = 60 * 24 * 60 * 60 // 60 Days

func main() {
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg.jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	log.Print(cfg.jwtSecret)
	cfg.db, _ = database.NewDB("database.json", *dbg)

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
	mux.HandleFunc("POST /api/login", postLoginHandler)
	mux.HandleFunc("POST /api/refresh", refreshTokenHandler)
	mux.HandleFunc("POST /api/revoke", revokeTokenHandler)

	mux.HandleFunc("PUT /api/users", putUsersHandler)

	loggingHandler := loggingMiddleware(mux)
	const port = "8080"
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
		log.Printf("Visited page: %v - %s\n", r.Method, r.URL.Path)
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

	chirp, err := cfg.db.CreateChirp(removeProfanity(c.Body))
	if err != nil {
		log.Print("failed to create chirp", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, chirp)
}

func getChirpHandler(w http.ResponseWriter, r *http.Request) {
	var chirps []database.Chirp

	chirps, err := cfg.db.GetChirps()
	if err != nil {
		log.Print("failure getting chirps", err)
	}

	id, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		respondWithJSON(w, http.StatusOK, chirps)
		return
	}

	c, err := cfg.db.GetChirp(id)
	if err != nil {
		respondWithError(w, 404, "chirp not found")
		return
	}
	respondWithJSON(w, http.StatusOK, c)
}

func postUsersHandler(w http.ResponseWriter, r *http.Request) {
	u := &struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := decodeJSON(r, u)
	if err != nil {
		log.Println("Error decoding params: ", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	user, err := cfg.db.CreateUser(u.Email, u.Password)
	if err != nil {
		log.Print("failed to create user", err)
		return
	}

	// remove password field from response
	response := struct {
		Email string `json:"email"`
		ID    int    `json:"id"`
	}{
		Email: user.Email,
		ID:    user.ID,
	}

	respondWithJSON(w, 201, response)
}

func postLoginHandler(w http.ResponseWriter, r *http.Request) {
	u := &struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Expires  int    `json:"expires_in_seconds,omitempty"`
	}{}

	err := decodeJSON(r, u)
	if err != nil {
		log.Println("Error decoding params: ", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	user, err := cfg.db.GetUser(u.Email)
	if err != nil {
		log.Print("failed to login", err)
		respondWithError(w, 401, "invalid user")
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(u.Password))
	if err != nil {
		// passwords do not match
		respondWithError(w, 401, "invalid password")
		return
	}

	if u.Expires < 1 || u.Expires >= accessTokenInterval {
		u.Expires = accessTokenInterval
	}

	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(u.Expires) * time.Second)),
		Subject:   strconv.Itoa(user.ID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(cfg.jwtSecret)
	if err != nil {
		log.Print("signing failed!", err)
		respondWithError(w, http.StatusInternalServerError, "signing failed")

	}

	// remove password field from response
	response := struct {
		Email        string `json:"email"`
		ID           int    `json:"id"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}{
		Email: user.Email,
		ID:    user.ID,
		Token: ss,
		RefreshToken: "CHANGE", // CHANGE THIS TO CORRECT 
	}

	respondWithJSON(w, 200, response)
}

func putUsersHandler(w http.ResponseWriter, r *http.Request) {
	// get token and verify auth

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		log.Println("Error: token not provided")
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	log.Println("Token String:", tokenString)

	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		return cfg.jwtSecret, nil
	})
	if err != nil {
		log.Printf("failed to parse token: %v", err)
		respondWithError(w, 401, "invalid token - parse error")
		return
	}

	if !token.Valid {
		respondWithError(w, 401, "invalid token")
		return
	}

	id, err := token.Claims.GetSubject()
	if err != nil {
		respondWithError(w, 401, "no user id provided in jwt")
		return
	}

	u := &struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err = decodeJSON(r, u)
	if err != nil {
		log.Println("Error decoding params: ", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	id2, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, 401, "invalid userID")
	}
	user := database.User{
		Email:    u.Email,
		Password: []byte(u.Password),
		ID:       id2,
	}

	user, err = cfg.db.UpdateUser(id2, user)
	// remove password field from response
	response := struct {
		Email string `json:"email"`
		ID    int    `json:"id"`
	}{
		Email: user.Email,
		ID:    user.ID,
	}

	respondWithJSON(w, 200, response)
}

func refreshTokenHandler(w http.ResponseWriter, r *http.Request) {

}

func revokeTokenHandler(w http.ResponseWriter, r *http.Request) {

}