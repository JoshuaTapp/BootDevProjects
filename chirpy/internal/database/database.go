package database

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	Password     []byte `json:"password"`
	RefreshToken string `json:"refresh_token"`
}

type Chirp struct {
	Body string `json:"body"`
	ID   int    `json:"id"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string, debug bool) (*DB, error) {
	db := DB{
		path: path,
		mux:  new(sync.RWMutex),
	}
	log.Printf("creating database: %v", db)
	err := db.ensureDB()
	if err != nil {
		return nil, err
	}
	if debug {
		// if debug, overwrite the file!
		initDB := new(DBStructure)
		initDB.Chirps = make(map[int]Chirp)
		initDB.Users = make(map[int]User)
		err = db.writeDB(*initDB)
		return &db, err
	}

	return &db, err
}

func (db *DB) CreateUser(email, password string) (u User, err error) {
	dbs, err := db.loadDB()
	if err != nil {
		return
	}
	newID := len(dbs.Users) + 1
	hashPW, err := bcrypt.GenerateFromPassword([]byte(password), 10) // default cost
	if err != nil {
		log.Fatal("password hashing failed!", password, hashPW)
	}
	u = User{
		Email:    email,
		ID:       newID,
		Password: hashPW,
	}

	dbs.Users[newID] = u
	err = db.writeDB(dbs)
	if err != nil {
		return User{}, err
	}

	return
}

func (db *DB) GetUser(email string) (u User, err error) {
	dbs, err := db.loadDB()
	if err != nil {
		return u, err
	}

	// this is dumb implementation, i should use emails instead for key...
	for _, v := range dbs.Users {
		if v.Email == email {
			u = v
			break
		}
	}
	if u.Email == "" {
		err = errors.New("user not found")
	}

	return
}

func (db *DB) GetUserById(id int) (u User, err error) {
	dbs, err := db.loadDB()
	if err != nil {
		return u, err
	}

	u, ok := dbs.Users[id]
	if !ok {
		return User{}, errors.New("user not found")
	}

	return u, nil
}

func (db *DB) UpdateUser(id int, user User) (User, error) {
	userOld, err := db.GetUserById(id)
	if err != nil {
		return userOld, err
	}

	if user.Email != "" {
		userOld.Email = user.Email
	}
	if len(user.Password) > 0 {
		hashPW, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10) // default cost
		if err != nil {
			log.Fatal("password hashing failed!", user.Password, hashPW)
		}
		user.Password = hashPW
	}
	user.ID = id

	dbs, err := db.loadDB()
	if err != nil {
		return userOld, err
	}

	dbs.Users[id] = user
	err = db.writeDB(dbs)
	if err != nil {
		return userOld, err
	}

	return user, nil
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
		log.Printf("failed to load during getChirps: %v", dbs)
		return
	}

	for _, chirp := range dbs.Chirps {
		chirps = append(chirps, chirp)
	}

	log.Println("GOT CHIRPS:", chirps, err)

	return
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirp(ID int) (chirp Chirp, err error) {
	dbs, err := db.loadDB()
	if err != nil {
		log.Printf("failed to load during getChirps: %v", dbs)
		return
	}

	if _, ok := dbs.Chirps[ID]; !ok {
		return chirp, errors.New("chirp does not exist")
	}

	log.Printf("requested ID: %v |\tChirp: %v", ID, dbs.Chirps[ID])
	return dbs.Chirps[ID], nil

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

func (db *DB) refreshUserRefreshToken(userID int) error {
	// generate new random 32 byte token
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Print("error refreshing token: ", err)
		return err
	}
	token := hex.EncodeToString(b)

	// assign new token to user
	dbs, err := db.loadDB()
	if err != nil {
		return err
	}

	if u, ok := dbs.Users[userID]; ok {
		u.RefreshToken = token
		dbs.Users[userID] = u
		return db.writeDB(dbs)
	} 
	
	return errors.New("can not refresh token because user does not exist")
	
}
