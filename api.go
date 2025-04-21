package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type SearchShowResult struct {
	Page    int
	Results []Show
}

type Show struct {
	Id           int
	OriginalName string `json:"original_name"`
}

func SearchShow(showName string) SearchShowResult {
	showName = strings.TrimSpace(showName)
	showName = strings.ReplaceAll(showName, " ", "%20")
	url := fmt.Sprintf("https://api.themoviedb.org/3/search/tv?query=%s&include_adult=false&language=en-US&page=1", showName)

	tmdbApiKey := GetDotenvValue("TMDB_API_KEY")

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+tmdbApiKey)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	var data SearchShowResult
	err := json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		return SearchShowResult{}
	}

	return data
}

func GetLastSeason(showID string) int {
	url := fmt.Sprintf("https://api.themoviedb.org/3/tv/%s?language=en-US", showID)

	tmdbApiKey := GetDotenvValue("TMDB_API_KEY")

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+tmdbApiKey)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	var data map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		return 0
	}

	lastSeason := int(data["number_of_seasons"].(float64))

	fmt.Println("Last season:", lastSeason)

	return lastSeason
}
