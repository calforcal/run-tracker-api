package storage

type (
	User struct {
		ID                  int
		UUID                string
		Name                string
		Username            string
		StravaID            int64
		StravaAccessToken   string
		StravaRefreshToken  string
		StravaExpiresAt     int
		SpotifyID           *string
		SpotifyAccessToken  *string
		SpotifyRefreshToken *string
		SpotifyExpiresAt    *int64
		CreatedAt           string
		UpdatedAt           string
	}

	WebhookSubscription struct {
		ID          int
		StravaID    int
		CallbackURL string
	}
)
