package repository

import (
	"log"
	"os"
)

type movieDB struct{}

func NewMovieDB() movieDB {
	return movieDB{}
}

func (db *movieDB) Log(record string) error {
	f, err := os.OpenFile("search.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	defer f.Close()

	logger := log.New(f, "movie_search: ", log.LstdFlags)
	logger.Println(record)

	return nil
}
