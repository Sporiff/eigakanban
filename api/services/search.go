package services

import (
	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"codeberg.org/sporiff/eigakanban/types"
	tmdb "github.com/cyruzin/golang-tmdb"
	"strconv"
)

type SearchService struct {
	q          *queries.Queries
	tmdbClient *tmdb.Client
}

func NewSearchService(q *queries.Queries, tmdbClient *tmdb.Client) *SearchService {
	return &SearchService{
		q:          q,
		tmdbClient: tmdbClient,
	}
}

// SearchMovie uses the TMDB API to search for movies
func (s *SearchService) SearchMovie(pagination *types.Pagination, q string) (*tmdb.SearchMovies, error) {
	var urlOptions = map[string]string{}
	parsedPage := strconv.Itoa(int(pagination.Page) + 1)
	urlOptions["page"] = parsedPage
	results, err := s.tmdbClient.GetSearchMovies(q, urlOptions)
	if err != nil {
		return nil, err
	}

	return results, nil
}
