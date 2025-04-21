package utils

import (
	"net/http"
)

func MakeApiRequest(url string) *http.Response {
	tmdbApiKey := GetDotenvValue("TMDB_API_KEY")

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+tmdbApiKey)

	res, _ := http.DefaultClient.Do(req)

	return res
}
