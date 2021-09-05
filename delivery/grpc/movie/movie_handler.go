package movie

import (
	context "context"
	"net/url"
	"strconv"
	"strings"

	model "github.com/zenkobert/sbtest-2/domain"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

var (
	incorrectImdbIDError = status.Error(codes.InvalidArgument, "incorrect IMDB ID")
	movieNotFoundError   = status.Error(codes.NotFound, "movie not found")
)

type movieServer struct {
	MovieUsecase model.MovieUsecase
}

func NewMovieServer(movieusecase model.MovieUsecase) SearchMovieServer {
	return &movieServer{
		MovieUsecase: movieusecase,
	}
}

func (serv *movieServer) SearchMovie(ctx context.Context, req *SearchMovieRequest) (resp *SearchMovieResponse, err error) {
	if req.Pagination <= 0 {
		req.Pagination = 1
	}

	req.Searchword = url.QueryEscape(req.Searchword)

	movieSearch, err := serv.MovieUsecase.SearchMovies(req.Searchword, uint32(req.Pagination))
	if err != nil {
		return resp, status.Error(codes.Internal, err.Error())
	}

	if movieSearch.Error != "" {
		return resp, movieNotFoundError
	}

	return serv.convertMovieSearchToRPCResponse(movieSearch), nil
}

func (serv *movieServer) GetMovieDetail(ctx context.Context, req *GetMovieDetailRequest) (resp *GetMovieDetailResponse, err error) {
	err = validateImdbID(req.Id)
	if err != nil {
		return resp, err
	}

	detail, err := serv.MovieUsecase.GetMovieDetailByID(req.Id)
	if err != nil {
		return resp, status.Error(codes.Internal, err.Error())
	}

	if detail.Error != "" {
		return resp, movieNotFoundError
	}

	return serv.convertMovieDetailToRPCResponse(detail), nil
}

func validateImdbID(id string) error {
	prefixIdx := strings.Index(id, "tt")
	if prefixIdx < 0 {
		return incorrectImdbIDError
	}

	id = strings.ReplaceAll(id, "tt", "")
	_, err := strconv.Atoi(id)
	if err != nil {
		return incorrectImdbIDError
	}

	return nil
}

func (serv *movieServer) convertMovieSearchToRPCResponse(m *model.MovieSearch) (r *SearchMovieResponse) {
	r = &SearchMovieResponse{}

	for _, s := range m.Search {
		r.Results = append(r.Results, &Search{
			Title:  s.Title,
			Year:   s.Type,
			ImdbId: s.ImdbID,
			Poster: s.Poster,
		})
	}

	r.Total = m.TotalResults

	return r
}

func (serv *movieServer) convertMovieDetailToRPCResponse(m *model.MovieDetail) (r *GetMovieDetailResponse) {
	r = &GetMovieDetailResponse{
		Title:      m.Title,
		Year:       m.Year,
		Rated:      m.Rated,
		Released:   m.Released,
		Runtime:    m.Runtime,
		Genre:      m.Genre,
		Director:   m.Director,
		Writer:     m.Writer,
		Actors:     m.Actors,
		Plot:       m.Plot,
		Language:   m.Language,
		Country:    m.Country,
		Awards:     m.Awards,
		Poster:     m.Poster,
		Metascore:  m.Metascore,
		ImdbRating: m.ImdbRating,
		ImdbVotes:  m.ImdbVotes,
		ImdbId:     m.ImdbID,
		Type:       m.Type,
		Dvd:        m.DVD,
		BoxOffice:  m.BoxOffice,
		Production: m.Production,
		Website:    m.Production,
	}

	for _, rating := range m.Ratings {
		r.Ratings = append(r.Ratings, &Rating{
			Source: rating.Source,
			Value:  rating.Value,
		})
	}

	return r
}
