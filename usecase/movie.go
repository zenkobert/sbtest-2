package usecase

import (
	model "github.com/zenkobert/sbtest-2/domain"
)

type movieUsecase struct {
	MovieRepo model.MovieRepository
}

func NewMovieUsecase(movieRepo model.MovieRepository) model.MovieUsecase {
	return &movieUsecase{
		MovieRepo: movieRepo,
	}
}

func (usecase *movieUsecase) SearchMovies(title string, page uint32) (result *model.MovieSearch, err error) {
	return usecase.MovieRepo.SearchMovies(title, page)
}

func (usecase *movieUsecase) GetMovieDetailByID(id string) (detail *model.MovieDetail, err error) {
	return usecase.MovieRepo.GetMovieDetailByID(id)
}
