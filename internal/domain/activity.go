package domain

import "time"

// ActivityAthlete is the minimal athlete reference embedded in an Activity.
type ActivityAthlete struct {
	ID            int64
	ResourceState int
}

// ActivityMap is the summary map embedded in an Activity.
type ActivityMap struct {
	ID              string
	SummaryPolyline *string
	ResourceState   int
}

// PolylineMap is the full map embedded in a DetailedActivity.
type PolylineMap struct {
	ActivityMap
	Polyline string
}

// Activity is a summary Strava activity, as returned from a list of activities.
type Activity struct {
	ResourceState        int
	Athlete              ActivityAthlete
	Name                 string
	Distance             float64
	MovingTime           int
	ElapsedTime          int
	TotalElevationGain   float64
	Type                 string
	SportType            string
	WorkoutType          *int
	ID                   int64
	ExternalID           string
	UploadID             int64
	StartDate            time.Time
	StartDateLocal       time.Time
	Timezone             string
	UTCOffset            float64
	StartLatLng          *[]float64
	EndLatLng            *[]float64
	LocationCity         *string
	LocationState        *string
	LocationCountry      *string
	AchievementCount     int
	KudosCount           int
	CommentCount         int
	AthleteCount         int
	PhotoCount           int
	Map                  ActivityMap
	Trainer              bool
	Commute              bool
	Manual               bool
	Private              bool
	Flagged              bool
	GearID               *string
	FromAcceptedTag      bool
	AverageSpeed         float64
	MaxSpeed             float64
	AverageCadence       *float64
	AverageWatts         *float64
	WeightedAverageWatts *float64
	Kilojoules           *float64
	DeviceWatts          *bool
	HasHeartrate         *bool
	AverageHeartrate     *float64
	MaxHeartrate         *float64
	MaxWatts             *float64
	PRCount              int
	TotalPhotoCount      int
	HasKudoed            bool
	SufferScore          *float64
}

// DetailedActivity is a single Strava activity with full detail.
type DetailedActivity struct {
	ID                 int64
	ExternalID         string
	UploadID           int64
	Name               string
	Distance           float64
	MovingTime         int
	ElapsedTime        int
	TotalElevationGain float64
	ElevHigh           float64
	ElevLow            float64
	Type               string
	SportType          string
	StartDate          time.Time
	StartDateLocal     time.Time
	Timezone           string
	StartLatLng        []float64
	EndLatLng          []float64
	AchievementCount   int
	KudosCount         int
	CommentCount       int
	AthleteCount       int
	PhotoCount         int
	TotalPhotoCount    int
	Map                PolylineMap
	Trainer            bool
	Commute            bool
	Manual             bool
	Private            bool
	Flagged            bool
	WorkoutType        int
	UploadIDStr        string
	AverageSpeed       float64
	MaxSpeed           float64
	HasKudoed          bool
	HideFromHome       bool
	GearID             string
	Kilojoules         float64
	AverageWatts       float64
	DeviceWatts        bool
	MaxWatts           int
	WeightedAvgWatts   int
	Description        string
	Calories           float64
	DeviceName         string
	EmbedToken         string
}

// ActivityStream is a single stream (e.g. heartrate, latlng) of activity data.
type ActivityStream struct {
	Type         string
	Data         []interface{}
	SeriesType   string
	OriginalSize int
	Resolution   string
}
