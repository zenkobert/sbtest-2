package repository

import (
	"log"
	"os"
	"sync"
)

type movieDB struct {
	mutex    *sync.Mutex
	fileName string
}

func NewMovieDB(fileName string) movieDB {
	return movieDB{fileName: fileName, mutex: &sync.Mutex{}}
}

func (db *movieDB) Log(record string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	f, err := os.OpenFile(db.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	defer f.Close()

	logger := log.New(f, "movie_search: ", log.LstdFlags)
	logger.Println(record)

	return nil
}
