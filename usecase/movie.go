package usecase

import (
	"github.com/zenkobert/sbtest-2/common"
	model "github.com/zenkobert/sbtest-2/domain"
)

type movieUsecase struct {
	MovieRepo model.MovieRepository
	MovieDB   common.DummyDB
}

func NewMovieUsecase(movieRepo model.MovieRepository, movieDB common.DummyDB) model.MovieUsecase {
	return &movieUsecase{
		MovieRepo: movieRepo,
		MovieDB:   movieDB,
	}
}

func (usecase *movieUsecase) SearchMovies(title string, page uint32) (result *model.MovieSearch, err error) {
	return usecase.MovieRepo.SearchMovies(title, page)
}

func (usecase *movieUsecase) GetMovieDetailByID(id string) (detail *model.MovieDetail, err error) {
	return usecase.MovieRepo.GetMovieDetailByID(id)
}

func (usecase *movieUsecase) LogToDB(record string) error {
	return usecase.MovieDB.Log(record)
}
