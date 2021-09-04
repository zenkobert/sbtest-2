package movie

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	model "github.com/zenkobert/sbtest-2/domain"

	"github.com/zenkobert/sbtest-2/common"
)

var host string = "http://www.omdbapi.com"

type movieRepo struct {
	Client common.HTTPClient
	apiKey string
}

func NewMovieRepo(apiKey string) model.MovieRepository {
	return &movieRepo{
		Client: &http.Client{},
		apiKey: apiKey,
	}
}

func (repo *movieRepo) SearchMovies(title string, page uint32) (result model.MovieSearch, err error) {
	url := fmt.Sprintf("%s/?apikey=%s&s=%s&page=%d", host, repo.apiKey, title, page)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return result, err
	}

	resp, err := repo.Client.Do(req)
	if err != nil || resp == nil {
		return result, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (repo *movieRepo) GetMovieDetailByID(id string) (detail model.MovieDetail, err error) {
	url := fmt.Sprintf("%s/?apikey=%s&i=%s", host, repo.apiKey, id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return detail, err
	}

	resp, err := repo.Client.Do(req)
	if err != nil || resp == nil {
		return detail, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return detail, err
	}

	err = json.Unmarshal(body, &detail)
	if err != nil {
		return detail, err
	}

	return detail, nil
}