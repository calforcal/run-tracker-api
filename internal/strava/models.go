package strava

type (
	Bike struct {
		ID            string  `json:"id"`
		Primary       bool    `json:"primary"`
		Name          string  `json:"name"`
		ResourceState int     `json:"resource_state"`
		Distance      float64 `json:"distance"`
	}

	Shoe struct {
		ID            string  `json:"id"`
		Primary       bool    `json:"primary"`
		Name          string  `json:"name"`
		ResourceState int     `json:"resource_state"`
		Distance      float64 `json:"distance"`
	}

	Athlete struct {
		ID                    int64         `json:"id"`
		Username              string        `json:"username"`
		ResourceState         int           `json:"resource_state"`
		Firstname             string        `json:"firstname"`
		Lastname              string        `json:"lastname"`
		City                  string        `json:"city"`
		State                 string        `json:"state"`
		Country               string        `json:"country"`
		Sex                   string        `json:"sex"`
		Premium               bool          `json:"premium"`
		CreatedAt             string        `json:"created_at"`
		UpdatedAt             string        `json:"updated_at"`
		BadgeTypeID           int           `json:"badge_type_id"`
		ProfileMedium         string        `json:"profile_medium"`
		Profile               string        `json:"profile"`
		Friend                *bool         `json:"friend"`
		Follower              *bool         `json:"follower"`
		FollowerCount         int           `json:"follower_count"`
		FriendCount           int           `json:"friend_count"`
		MutualFriendCount     int           `json:"mutual_friend_count"`
		AthleteType           int           `json:"athlete_type"`
		DatePreference        string        `json:"date_preference"`
		MeasurementPreference string        `json:"measurement_preference"`
		Clubs                 []interface{} `json:"clubs"`
		FTP                   *float64      `json:"ftp"`
		Weight                float64       `json:"weight"`
		Bikes                 []Bike        `json:"bikes"`
		Shoes                 []Shoe        `json:"shoes"`
		IsSpotifyConnected    *bool         `json:"is_spotify_connected"`
	}

	Activities struct {
	}
)

type ActivityAthlete struct {
	ID            int64 `json:"id"`
	ResourceState int   `json:"resource_state"`
}

type ActivityMap struct {
	ID              string  `json:"id"`
	SummaryPolyline *string `json:"summary_polyline"`
	ResourceState   int     `json:"resource_state"`
}

type Activity struct {
	ResourceState        int             `json:"resource_state"`
	Athlete              ActivityAthlete `json:"athlete"`
	Name                 string          `json:"name"`
	Distance             float64         `json:"distance"`
	MovingTime           int             `json:"moving_time"`
	ElapsedTime          int             `json:"elapsed_time"`
	TotalElevationGain   float64         `json:"total_elevation_gain"`
	Type                 string          `json:"type"`
	SportType            string          `json:"sport_type"`
	WorkoutType          *int            `json:"workout_type"`
	ID                   int64           `json:"id"`
	ExternalID           string          `json:"external_id"`
	UploadID             int64           `json:"upload_id"`
	StartDate            string          `json:"start_date"`
	StartDateLocal       string          `json:"start_date_local"`
	Timezone             string          `json:"timezone"`
	UTCOffset            float64         `json:"utc_offset"`
	StartLatLng          *[]float64      `json:"start_latlng"`
	EndLatLng            *[]float64      `json:"end_latlng"`
	LocationCity         *string         `json:"location_city"`
	LocationState        *string         `json:"location_state"`
	LocationCountry      *string         `json:"location_country"`
	AchievementCount     int             `json:"achievement_count"`
	KudosCount           int             `json:"kudos_count"`
	CommentCount         int             `json:"comment_count"`
	AthleteCount         int             `json:"athlete_count"`
	PhotoCount           int             `json:"photo_count"`
	Map                  ActivityMap     `json:"map"`
	Trainer              bool            `json:"trainer"`
	Commute              bool            `json:"commute"`
	Manual               bool            `json:"manual"`
	Private              bool            `json:"private"`
	Flagged              bool            `json:"flagged"`
	GearID               *string         `json:"gear_id"`
	FromAcceptedTag      bool            `json:"from_accepted_tag"`
	AverageSpeed         float64         `json:"average_speed"`
	MaxSpeed             float64         `json:"max_speed"`
	AverageCadence       *float64        `json:"average_cadence"`
	AverageWatts         *float64        `json:"average_watts"`
	WeightedAverageWatts *float64        `json:"weighted_average_watts"`
	Kilojoules           *float64        `json:"kilojoules"`
	DeviceWatts          *bool           `json:"device_watts"`
	HasHeartrate         *bool           `json:"has_heartrate"`
	AverageHeartrate     *float64        `json:"average_heartrate"`
	MaxHeartrate         *float64        `json:"max_heartrate"`
	MaxWatts             *float64        `json:"max_watts"`
	PRCount              int             `json:"pr_count"`
	TotalPhotoCount      int             `json:"total_photo_count"`
	HasKudoed            bool            `json:"has_kudoed"`
	SufferScore          *float64        `json:"suffer_score"`
}

type PolylineMap struct {
	ActivityMap
	Polyline string `json:"polyline"`
}

type DetailedActivity struct {
	ID                 int64       `json:"id"`
	ExternalID         string      `json:"external_id"`
	UploadID           int64       `json:"upload_id"`
	Name               string      `json:"name"`
	Distance           float64     `json:"distance"`
	MovingTime         int         `json:"moving_time"`
	ElapsedTime        int         `json:"elapsed_time"`
	TotalElevationGain float64     `json:"total_elevation_gain"`
	ElevHigh           float64     `json:"elev_high"`
	ElevLow            float64     `json:"elev_low"`
	Type               string      `json:"type"`       // deprecated, prefer SportType
	SportType          string      `json:"sport_type"` // could be enum if you want
	StartDate          string      `json:"start_date"`
	StartDateLocal     string      `json:"start_date_local"`
	Timezone           string      `json:"timezone"`
	StartLatLng        []float64   `json:"start_latlng"`
	EndLatLng          []float64   `json:"end_latlng"`
	AchievementCount   int         `json:"achievement_count"`
	KudosCount         int         `json:"kudos_count"`
	CommentCount       int         `json:"comment_count"`
	AthleteCount       int         `json:"athlete_count"`
	PhotoCount         int         `json:"photo_count"`
	TotalPhotoCount    int         `json:"total_photo_count"`
	Map                PolylineMap `json:"map"`
	Trainer            bool        `json:"trainer"`
	Commute            bool        `json:"commute"`
	Manual             bool        `json:"manual"`
	Private            bool        `json:"private"`
	Flagged            bool        `json:"flagged"`
	WorkoutType        int         `json:"workout_type"`
	UploadIDStr        string      `json:"upload_id_str"`
	AverageSpeed       float64     `json:"average_speed"`
	MaxSpeed           float64     `json:"max_speed"`
	HasKudoed          bool        `json:"has_kudoed"`
	HideFromHome       bool        `json:"hide_from_home"`
	GearID             string      `json:"gear_id"`
	Kilojoules         float64     `json:"kilojoules"`
	AverageWatts       float64     `json:"average_watts"`
	DeviceWatts        bool        `json:"device_watts"`
	MaxWatts           int         `json:"max_watts"`
	WeightedAvgWatts   int         `json:"weighted_average_watts"`
	Description        string
	Calories           float64 `json:"calories"`
	DeviceName         string  `json:"device_name"`
	EmbedToken         string  `json:"embed_token"`
}
