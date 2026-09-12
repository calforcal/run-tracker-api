package dto

import (
	"time"

	"run-tracker-api/internal/domain"
)

// These response types mirror the exact JSON shape the frontend has always
// received (originally the raw Strava wire structs were serialized
// directly), now built explicitly from domain.Athlete/Activity so the API
// layer - not the provider - owns the outbound contract.

type BikeResponse struct {
	ID            string  `json:"id"`
	Primary       bool    `json:"primary"`
	Name          string  `json:"name"`
	ResourceState int     `json:"resource_state"`
	Distance      float64 `json:"distance"`
}

type ShoeResponse struct {
	ID            string  `json:"id"`
	Primary       bool    `json:"primary"`
	Name          string  `json:"name"`
	ResourceState int     `json:"resource_state"`
	Distance      float64 `json:"distance"`
}

type AthleteResponse struct {
	ID                    int64          `json:"id"`
	Username              string         `json:"username"`
	ResourceState         int            `json:"resource_state"`
	Firstname             string         `json:"firstname"`
	Lastname              string         `json:"lastname"`
	City                  string         `json:"city"`
	State                 string         `json:"state"`
	Country               string         `json:"country"`
	Sex                   string         `json:"sex"`
	Premium               bool           `json:"premium"`
	CreatedAt             string         `json:"created_at"`
	UpdatedAt             string         `json:"updated_at"`
	BadgeTypeID           int            `json:"badge_type_id"`
	ProfileMedium         string         `json:"profile_medium"`
	Profile               string         `json:"profile"`
	Friend                *bool          `json:"friend"`
	Follower              *bool          `json:"follower"`
	FollowerCount         int            `json:"follower_count"`
	FriendCount           int            `json:"friend_count"`
	MutualFriendCount     int            `json:"mutual_friend_count"`
	AthleteType           int            `json:"athlete_type"`
	DatePreference        string         `json:"date_preference"`
	MeasurementPreference string         `json:"measurement_preference"`
	FTP                   *float64       `json:"ftp"`
	Weight                float64        `json:"weight"`
	Bikes                 []BikeResponse `json:"bikes"`
	Shoes                 []ShoeResponse `json:"shoes"`
	IsSpotifyConnected    bool           `json:"is_spotify_connected"`
}

// AthleteFromDomain builds an AthleteResponse. isSpotifyConnected is
// computed by the caller (the handler already has the domain.User) rather
// than living on domain.Athlete, which is a pure Strava concept.
func AthleteFromDomain(a domain.Athlete, isSpotifyConnected bool) AthleteResponse {
	bikes := make([]BikeResponse, len(a.Bikes))
	for i, b := range a.Bikes {
		bikes[i] = BikeResponse{ID: b.ID, Primary: b.Primary, Name: b.Name, ResourceState: b.ResourceState, Distance: b.Distance}
	}
	shoes := make([]ShoeResponse, len(a.Shoes))
	for i, sh := range a.Shoes {
		shoes[i] = ShoeResponse{ID: sh.ID, Primary: sh.Primary, Name: sh.Name, ResourceState: sh.ResourceState, Distance: sh.Distance}
	}

	return AthleteResponse{
		ID:                    a.ID,
		Username:              a.Username,
		ResourceState:         a.ResourceState,
		Firstname:             a.Firstname,
		Lastname:              a.Lastname,
		City:                  a.City,
		State:                 a.State,
		Country:               a.Country,
		Sex:                   a.Sex,
		Premium:               a.Premium,
		CreatedAt:             formatTime(a.CreatedAt),
		UpdatedAt:             formatTime(a.UpdatedAt),
		BadgeTypeID:           a.BadgeTypeID,
		ProfileMedium:         a.ProfileMedium,
		Profile:               a.Profile,
		Friend:                a.Friend,
		Follower:              a.Follower,
		FollowerCount:         a.FollowerCount,
		FriendCount:           a.FriendCount,
		MutualFriendCount:     a.MutualFriendCount,
		AthleteType:           a.AthleteType,
		DatePreference:        a.DatePreference,
		MeasurementPreference: a.MeasurementPreference,
		FTP:                   a.FTP,
		Weight:                a.Weight,
		Bikes:                 bikes,
		Shoes:                 shoes,
		IsSpotifyConnected:    isSpotifyConnected,
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

type ActivityAthleteResponse struct {
	ID            int64 `json:"id"`
	ResourceState int   `json:"resource_state"`
}

type ActivityMapResponse struct {
	ID              string  `json:"id"`
	SummaryPolyline *string `json:"summary_polyline"`
	ResourceState   int     `json:"resource_state"`
}

type PolylineMapResponse struct {
	ActivityMapResponse
	Polyline string `json:"polyline"`
}

type ActivityResponse struct {
	ResourceState        int                     `json:"resource_state"`
	Athlete              ActivityAthleteResponse `json:"athlete"`
	Name                 string                  `json:"name"`
	Distance             float64                 `json:"distance"`
	MovingTime           int                     `json:"moving_time"`
	ElapsedTime          int                     `json:"elapsed_time"`
	TotalElevationGain   float64                 `json:"total_elevation_gain"`
	Type                 string                  `json:"type"`
	SportType            string                  `json:"sport_type"`
	WorkoutType          *int                    `json:"workout_type"`
	ID                   int64                   `json:"id"`
	ExternalID           string                  `json:"external_id"`
	UploadID             int64                   `json:"upload_id"`
	StartDate            string                  `json:"start_date"`
	StartDateLocal       string                  `json:"start_date_local"`
	Timezone             string                  `json:"timezone"`
	UTCOffset            float64                 `json:"utc_offset"`
	StartLatLng          *[]float64              `json:"start_latlng"`
	EndLatLng            *[]float64              `json:"end_latlng"`
	LocationCity         *string                 `json:"location_city"`
	LocationState        *string                 `json:"location_state"`
	LocationCountry      *string                 `json:"location_country"`
	AchievementCount     int                     `json:"achievement_count"`
	KudosCount           int                     `json:"kudos_count"`
	CommentCount         int                     `json:"comment_count"`
	AthleteCount         int                     `json:"athlete_count"`
	PhotoCount           int                     `json:"photo_count"`
	Map                  ActivityMapResponse     `json:"map"`
	Trainer              bool                    `json:"trainer"`
	Commute              bool                    `json:"commute"`
	Manual               bool                    `json:"manual"`
	Private              bool                    `json:"private"`
	Flagged              bool                    `json:"flagged"`
	GearID               *string                 `json:"gear_id"`
	FromAcceptedTag      bool                    `json:"from_accepted_tag"`
	AverageSpeed         float64                 `json:"average_speed"`
	MaxSpeed             float64                 `json:"max_speed"`
	AverageCadence       *float64                `json:"average_cadence"`
	AverageWatts         *float64                `json:"average_watts"`
	WeightedAverageWatts *float64                `json:"weighted_average_watts"`
	Kilojoules           *float64                `json:"kilojoules"`
	DeviceWatts          *bool                   `json:"device_watts"`
	HasHeartrate         *bool                   `json:"has_heartrate"`
	AverageHeartrate     *float64                `json:"average_heartrate"`
	MaxHeartrate         *float64                `json:"max_heartrate"`
	MaxWatts             *float64                `json:"max_watts"`
	PRCount              int                     `json:"pr_count"`
	TotalPhotoCount      int                     `json:"total_photo_count"`
	HasKudoed            bool                    `json:"has_kudoed"`
	SufferScore          *float64                `json:"suffer_score"`
}

func ActivityFromDomain(a domain.Activity) ActivityResponse {
	return ActivityResponse{
		ResourceState:        a.ResourceState,
		Athlete:              ActivityAthleteResponse{ID: a.Athlete.ID, ResourceState: a.Athlete.ResourceState},
		Name:                 a.Name,
		Distance:             a.Distance,
		MovingTime:           a.MovingTime,
		ElapsedTime:          a.ElapsedTime,
		TotalElevationGain:   a.TotalElevationGain,
		Type:                 a.Type,
		SportType:            a.SportType,
		WorkoutType:          a.WorkoutType,
		ID:                   a.ID,
		ExternalID:           a.ExternalID,
		UploadID:             a.UploadID,
		StartDate:            formatTime(a.StartDate),
		StartDateLocal:       formatTime(a.StartDateLocal),
		Timezone:             a.Timezone,
		UTCOffset:            a.UTCOffset,
		StartLatLng:          a.StartLatLng,
		EndLatLng:            a.EndLatLng,
		LocationCity:         a.LocationCity,
		LocationState:        a.LocationState,
		LocationCountry:      a.LocationCountry,
		AchievementCount:     a.AchievementCount,
		KudosCount:           a.KudosCount,
		CommentCount:         a.CommentCount,
		AthleteCount:         a.AthleteCount,
		PhotoCount:           a.PhotoCount,
		Map:                  ActivityMapResponse{ID: a.Map.ID, SummaryPolyline: a.Map.SummaryPolyline, ResourceState: a.Map.ResourceState},
		Trainer:              a.Trainer,
		Commute:              a.Commute,
		Manual:               a.Manual,
		Private:              a.Private,
		Flagged:              a.Flagged,
		GearID:               a.GearID,
		FromAcceptedTag:      a.FromAcceptedTag,
		AverageSpeed:         a.AverageSpeed,
		MaxSpeed:             a.MaxSpeed,
		AverageCadence:       a.AverageCadence,
		AverageWatts:         a.AverageWatts,
		WeightedAverageWatts: a.WeightedAverageWatts,
		Kilojoules:           a.Kilojoules,
		DeviceWatts:          a.DeviceWatts,
		HasHeartrate:         a.HasHeartrate,
		AverageHeartrate:     a.AverageHeartrate,
		MaxHeartrate:         a.MaxHeartrate,
		MaxWatts:             a.MaxWatts,
		PRCount:              a.PRCount,
		TotalPhotoCount:      a.TotalPhotoCount,
		HasKudoed:            a.HasKudoed,
		SufferScore:          a.SufferScore,
	}
}

func ActivitiesFromDomain(activities []domain.Activity) []ActivityResponse {
	out := make([]ActivityResponse, len(activities))
	for i, a := range activities {
		out[i] = ActivityFromDomain(a)
	}
	return out
}

type DetailedActivityResponse struct {
	ID                 int64               `json:"id"`
	ExternalID         string              `json:"external_id"`
	UploadID           int64               `json:"upload_id"`
	Name               string              `json:"name"`
	Distance           float64             `json:"distance"`
	MovingTime         int                 `json:"moving_time"`
	ElapsedTime        int                 `json:"elapsed_time"`
	TotalElevationGain float64             `json:"total_elevation_gain"`
	ElevHigh           float64             `json:"elev_high"`
	ElevLow            float64             `json:"elev_low"`
	Type               string              `json:"type"`
	SportType          string              `json:"sport_type"`
	StartDate          string              `json:"start_date"`
	StartDateLocal     string              `json:"start_date_local"`
	Timezone           string              `json:"timezone"`
	StartLatLng        []float64           `json:"start_latlng"`
	EndLatLng          []float64           `json:"end_latlng"`
	AchievementCount   int                 `json:"achievement_count"`
	KudosCount         int                 `json:"kudos_count"`
	CommentCount       int                 `json:"comment_count"`
	AthleteCount       int                 `json:"athlete_count"`
	PhotoCount         int                 `json:"photo_count"`
	TotalPhotoCount    int                 `json:"total_photo_count"`
	Map                PolylineMapResponse `json:"map"`
	Trainer            bool                `json:"trainer"`
	Commute            bool                `json:"commute"`
	Manual             bool                `json:"manual"`
	Private            bool                `json:"private"`
	Flagged            bool                `json:"flagged"`
	WorkoutType        int                 `json:"workout_type"`
	UploadIDStr        string              `json:"upload_id_str"`
	AverageSpeed       float64             `json:"average_speed"`
	MaxSpeed           float64             `json:"max_speed"`
	HasKudoed          bool                `json:"has_kudoed"`
	HideFromHome       bool                `json:"hide_from_home"`
	GearID             string              `json:"gear_id"`
	Kilojoules         float64             `json:"kilojoules"`
	AverageWatts       float64             `json:"average_watts"`
	DeviceWatts        bool                `json:"device_watts"`
	MaxWatts           int                 `json:"max_watts"`
	WeightedAvgWatts   int                 `json:"weighted_average_watts"`
	// Description intentionally has no explicit json tag, matching the
	// exact (inconsistently-capitalized) key the frontend has always
	// received from this endpoint.
	Description string  `json:"Description"`
	Calories    float64 `json:"calories"`
	DeviceName  string  `json:"device_name"`
	EmbedToken  string  `json:"embed_token"`
}

func DetailedActivityFromDomain(a domain.DetailedActivity) DetailedActivityResponse {
	return DetailedActivityResponse{
		ID:                 a.ID,
		ExternalID:         a.ExternalID,
		UploadID:           a.UploadID,
		Name:               a.Name,
		Distance:           a.Distance,
		MovingTime:         a.MovingTime,
		ElapsedTime:        a.ElapsedTime,
		TotalElevationGain: a.TotalElevationGain,
		ElevHigh:           a.ElevHigh,
		ElevLow:            a.ElevLow,
		Type:               a.Type,
		SportType:          a.SportType,
		StartDate:          formatTime(a.StartDate),
		StartDateLocal:     formatTime(a.StartDateLocal),
		Timezone:           a.Timezone,
		StartLatLng:        a.StartLatLng,
		EndLatLng:          a.EndLatLng,
		AchievementCount:   a.AchievementCount,
		KudosCount:         a.KudosCount,
		CommentCount:       a.CommentCount,
		AthleteCount:       a.AthleteCount,
		PhotoCount:         a.PhotoCount,
		TotalPhotoCount:    a.TotalPhotoCount,
		Map: PolylineMapResponse{
			ActivityMapResponse: ActivityMapResponse{ID: a.Map.ID, SummaryPolyline: a.Map.SummaryPolyline, ResourceState: a.Map.ResourceState},
			Polyline:            a.Map.Polyline,
		},
		Trainer:          a.Trainer,
		Commute:          a.Commute,
		Manual:           a.Manual,
		Private:          a.Private,
		Flagged:          a.Flagged,
		WorkoutType:      a.WorkoutType,
		UploadIDStr:      a.UploadIDStr,
		AverageSpeed:     a.AverageSpeed,
		MaxSpeed:         a.MaxSpeed,
		HasKudoed:        a.HasKudoed,
		HideFromHome:     a.HideFromHome,
		GearID:           a.GearID,
		Kilojoules:       a.Kilojoules,
		AverageWatts:     a.AverageWatts,
		DeviceWatts:      a.DeviceWatts,
		MaxWatts:         a.MaxWatts,
		WeightedAvgWatts: a.WeightedAvgWatts,
		Description:      a.Description,
		Calories:         a.Calories,
		DeviceName:       a.DeviceName,
		EmbedToken:       a.EmbedToken,
	}
}

type ActivityStreamResponse struct {
	Type         string        `json:"type"`
	Data         []interface{} `json:"data"`
	SeriesType   string        `json:"series_type"`
	OriginalSize int           `json:"original_size"`
	Resolution   string        `json:"resolution"`
}

func ActivityStreamsFromDomain(streams []domain.ActivityStream) []ActivityStreamResponse {
	out := make([]ActivityStreamResponse, len(streams))
	for i, s := range streams {
		out[i] = ActivityStreamResponse{
			Type:         s.Type,
			Data:         s.Data,
			SeriesType:   s.SeriesType,
			OriginalSize: s.OriginalSize,
			Resolution:   s.Resolution,
		}
	}
	return out
}
