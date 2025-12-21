package spotify

type (
	SpotifyUser struct {
		DisplayName string `json:"display_name"`
		ID          string `json:"id"`
		Email       string `json:"email"`
	}

	ListeningHistory struct {
		Items []ListeningHistoryItem `json:"items"`
	}

	ListeningHistoryItem struct {
		Track    TrackInfo `json:"track"`
		PlayedAt string    `json:"played_at"`
	}

	TrackInfo struct {
		Album      AlbumInfo `json:"album"`
		Artists    []Artist  `json:"artists"`
		DurationMs int       `json:"duration_ms"`
		ID         string    `json:"id"`
		Name       string    `json:"name"`
		URI        string    `json:"uri"`
	}

	AlbumInfo struct {
		Name        string  `json:"name"`
		ReleaseDate string  `json:"release_date"`
		Images      []Image `json:"images"`
	}

	Artist struct {
		Name string `json:"name"`
	}

	Image struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	}
)
