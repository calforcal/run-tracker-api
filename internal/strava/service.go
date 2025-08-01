package strava

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type StravaService struct {
	client *http.Client
}

func New() *StravaService {
	return &StravaService{
		client: &http.Client{},
	}
}

func (s *StravaService) GetAthlete(c echo.Context) (Athlete, error) {
	req, err := http.NewRequest("GET", "https://www.strava.com/api/v3/athlete", nil)
	req.Header.Set("Authorization", "Bearer "+os.Getenv("STRAVA_ACCESS_TOKEN"))
	if err != nil {
		return Athlete{}, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return Athlete{}, err
	}
	defer resp.Body.Close()

	jsonDecoder := json.NewDecoder(resp.Body)
	var athlete Athlete
	err = jsonDecoder.Decode(&athlete)
	if err != nil {
		return Athlete{}, err
	}

	return athlete, nil
}
