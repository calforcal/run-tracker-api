package storage

type (
	User struct {
		ID                 int
		UUID               string
		Name               string
		Username           string
		StravaID           int64
		StravaAccessToken  string
		StravaRefreshToken string
		StravaExpiresAt    int
		CreatedAt          string
		UpdatedAt          string
	}
)
