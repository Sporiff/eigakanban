package services

import (
	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"codeberg.org/sporiff/eigakanban/types"
	tmdb "github.com/cyruzin/golang-tmdb"
	"net/http"
	"strconv"
)

type SearchService struct {
	q          *queries.Queries
	TMDBClient *tmdb.Client
}

func NewSearchService(q *queries.Queries, tmdbClient *tmdb.Client) *SearchService {
	return &SearchService{
		q:          q,
		TMDBClient: tmdbClient,
	}
}

// SearchMovie uses the TMDB API to search for movies
func (s *SearchService) SearchMovie(pagination *types.Pagination, q string) (*tmdb.SearchMovies, error) {
	var urlOptions = map[string]string{}
	parsedPage := strconv.Itoa(int(pagination.Page) + 1)
	urlOptions["page"] = parsedPage
	results, err := s.TMDBClient.GetSearchMovies(q, urlOptions)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "failed to fetch movies")
	}

	return results, nil
}
