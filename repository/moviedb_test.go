package repository

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMovieDB(t *testing.T) {
	t.Run("[NewMovieDB]", func(t *testing.T) {
		fileName := "search.log"
		assert.Equal(t, movieDB{fileName: fileName}, NewMovieDB(fileName))
	})
}

func TestLog(t *testing.T) {
	t.Run("[Log] assert file content", func(t *testing.T) {
		tempFileName := "tempFile.log"
		logContent := "test"

		movieDB := NewMovieDB(tempFileName)
		movieDB.Log(logContent)

		dat, err := os.ReadFile(tempFileName)
		defer os.Remove(tempFileName)
		if assert.Nil(t, err) {
			assert.Contains(t, string(dat), logContent)
		}
	})

	t.Run("[Log] return error if malformed filename", func(t *testing.T) {
		tempFileName := "tempF///ilelog"
		logContent := "test"

		movieDB := NewMovieDB(tempFileName)
		movieDB.Log(logContent)

		_, err := os.ReadFile(tempFileName)
		defer os.Remove(tempFileName)
		assert.Error(t, err)
	})
}
