package movie

import (
	context "context"

	model "github.com/zenkobert/sbtest-2/domain"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
	movieSearch, err := serv.MovieUsecase.SearchMovies(req.Searchword, uint32(req.Pagination))
	if err != nil {
		return resp, status.Error(codes.Internal, "Oops, something happened")
	}

	return serv.convertMovieSearchToRPCResponse(movieSearch), nil
}
func (serv *movieServer) GetMovieDetail(ctx context.Context, req *GetMovieDetailRequest) (resp *GetMovieDetailResponse, err error) {
	detail, err := serv.MovieUsecase.GetMovieDetailByID(req.Id)
	if err != nil {
		return resp, status.Error(codes.Internal, "Oops, something happened")
	}

	return serv.convertMovieDetailToRPCResponse(detail), nil
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
