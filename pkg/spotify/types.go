package spotify

// Wire-format types matching Spotify's public API JSON responses.

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type SpotifyUser struct {
	DisplayName string `json:"display_name"`
	ID          string `json:"id"`
	Email       string `json:"email"`
}

type ListeningHistory struct {
	Items []ListeningHistoryItem `json:"items"`
}

type ListeningHistoryItem struct {
	Track    TrackInfo `json:"track"`
	PlayedAt string    `json:"played_at"`
}

type TrackInfo struct {
	Album      AlbumInfo `json:"album"`
	Artists    []Artist  `json:"artists"`
	DurationMs int       `json:"duration_ms"`
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	URI        string    `json:"uri"`
}

type AlbumInfo struct {
	Name        string  `json:"name"`
	ReleaseDate string  `json:"release_date"`
	Images      []Image `json:"images"`
}

type Artist struct {
	Name string `json:"name"`
}

type Image struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}
