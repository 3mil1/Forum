package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
)

type Storage struct {
	config   *Config
	db       *sql.DB
	userRepo *UserRepo
	authRepo *AuthRepo
	postRepo *PostRepo
}

func New(config *Config) *Storage {
	return &Storage{
		config: config,
	}
}

func (storage *Storage) Open() error {
	db, err := sql.Open("sqlite3", "file:"+storage.config.DatabaseURI+"?_foreign_keys=on")

	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	storage.db = db
	log.Println("Connection to db successfully")
	return nil
}

func (storage *Storage) Close() {
	storage.db.Close()
}

func (storage *Storage) AddTables() {
	createSQL, err := ioutil.ReadFile("./createTables.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = storage.db.Exec(string(createSQL))
	if err != nil {
		log.Fatal(err)
	}
}

// User Public repo for user
func (storage *Storage) User() *UserRepo {
	if storage.userRepo != nil {
		return storage.userRepo
	}

	storage.userRepo = &UserRepo{
		storage: storage,
	}
	return storage.userRepo
}

// Auth Public repo
func (storage *Storage) Auth() *AuthRepo {
	if storage.authRepo != nil {
		return storage.authRepo
	}

	storage.authRepo = &AuthRepo{
		storage: storage,
	}
	return storage.authRepo
}

// Post Public repo
func (storage *Storage) Post() *PostRepo {
	if storage.postRepo != nil {
		return storage.postRepo
	}

	storage.postRepo = &PostRepo{
		storage: storage,
	}
	return storage.postRepo
}
