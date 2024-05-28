package database

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	Body string `json:"body"`
	ID   int    `json:"id"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := DB{
		path: path,
		mux:  new(sync.RWMutex),
	}
	log.Printf("creating database: %v", db)
	err := db.ensureDB()
	if err != nil {
		return nil, err
	}

	initDB := new(DBStructure)
	initDB.Chirps = make(map[int]Chirp)
	err = db.writeDB(*initDB)
	return &db, err
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (c Chirp, err error) {
	dbs, err := db.loadDB()
	if err != nil {
		return
	}
	newID := len(dbs.Chirps) + 1
	c = Chirp{
		Body: body,
		ID:   newID,
	}

	dbs.Chirps[newID] = c
	err = db.writeDB(dbs)
	if err != nil {
		return Chirp{}, err
	}

	return
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() (chirps []Chirp, err error) {
	dbs, err := db.loadDB()
	if err != nil {
		log.Printf("failed to load during getChirps: %v", dbs, err)
		return
	}

	for _, chirp := range dbs.Chirps {
		chirps = append(chirps, chirp)
	}

	log.Println("GOT CHIRPS:", chirps, err)

	return
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	dbFile, err := os.OpenFile(db.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Print("Failed ensureDB check!")
		return err
	}
	defer dbFile.Close()

	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (dbs DBStructure, err error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbFile, err := os.Open(db.path)
	if err != nil {
		return
	}
	defer dbFile.Close()

	byteVal, err := io.ReadAll(dbFile)
	if err != nil {
		return
	}

	err = json.Unmarshal(byteVal, &dbs)
	if err != nil {
		return
	}

	return dbs, nil

}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbs DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbFile, err := os.OpenFile(db.path, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer dbFile.Close()

	byteVal, err := json.Marshal(dbs)
	if err != nil {
		return err
	}

	_, err = dbFile.Write(byteVal)
	if err != nil {
		return err
	}

	return nil

}
