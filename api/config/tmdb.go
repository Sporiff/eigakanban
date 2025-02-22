package config

import (
	tmdb "github.com/cyruzin/golang-tmdb"
	"os"
)

func LoadTmdbConfig() (*tmdb.Client, error) {
	tmdbClient, err := tmdb.Init(os.Getenv("TMDB_API_KEY"))
	if err != nil {
		return nil, err
	}

	tmdbClient.SetClientAutoRetry()

	return tmdbClient, nil
}
