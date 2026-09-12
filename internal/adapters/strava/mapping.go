package strava

import (
	"time"

	"run-tracker-api/internal/domain"
	stravapkg "run-tracker-api/pkg/strava"
)

// parseTime parses one of Strava's RFC3339 date strings. Strava always
// returns well-formed dates for these fields; on the rare malformed
// response we fall back to the zero time rather than failing the whole
// request, matching how the old code silently passed the raw (possibly
// unparsed) string straight through.
func parseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func toDomainBikes(bikes []stravapkg.Bike) []domain.Bike {
	out := make([]domain.Bike, len(bikes))
	for i, b := range bikes {
		out[i] = domain.Bike{ID: b.ID, Name: b.Name, Primary: b.Primary, ResourceState: b.ResourceState, Distance: b.Distance}
	}
	return out
}

func toDomainShoes(shoes []stravapkg.Shoe) []domain.Shoe {
	out := make([]domain.Shoe, len(shoes))
	for i, s := range shoes {
		out[i] = domain.Shoe{ID: s.ID, Name: s.Name, Primary: s.Primary, ResourceState: s.ResourceState, Distance: s.Distance}
	}
	return out
}

func toDomainAthlete(a stravapkg.Athlete) domain.Athlete {
	return domain.Athlete{
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
		CreatedAt:             parseTime(a.CreatedAt),
		UpdatedAt:             parseTime(a.UpdatedAt),
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
		Bikes:                 toDomainBikes(a.Bikes),
		Shoes:                 toDomainShoes(a.Shoes),
	}
}

func toDomainActivityMap(m stravapkg.ActivityMap) domain.ActivityMap {
	return domain.ActivityMap{ID: m.ID, SummaryPolyline: m.SummaryPolyline, ResourceState: m.ResourceState}
}

func toDomainActivity(a stravapkg.Activity) domain.Activity {
	return domain.Activity{
		ResourceState:        a.ResourceState,
		Athlete:              domain.ActivityAthlete{ID: a.Athlete.ID, ResourceState: a.Athlete.ResourceState},
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
		StartDate:            parseTime(a.StartDate),
		StartDateLocal:       parseTime(a.StartDateLocal),
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
		Map:                  toDomainActivityMap(a.Map),
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

func toDomainActivities(activities []stravapkg.Activity) []domain.Activity {
	out := make([]domain.Activity, len(activities))
	for i, a := range activities {
		out[i] = toDomainActivity(a)
	}
	return out
}

func toDomainDetailedActivity(a stravapkg.DetailedActivity) domain.DetailedActivity {
	return domain.DetailedActivity{
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
		StartDate:          parseTime(a.StartDate),
		StartDateLocal:     parseTime(a.StartDateLocal),
		Timezone:           a.Timezone,
		StartLatLng:        a.StartLatLng,
		EndLatLng:          a.EndLatLng,
		AchievementCount:   a.AchievementCount,
		KudosCount:         a.KudosCount,
		CommentCount:       a.CommentCount,
		AthleteCount:       a.AthleteCount,
		PhotoCount:         a.PhotoCount,
		TotalPhotoCount:    a.TotalPhotoCount,
		Map: domain.PolylineMap{
			ActivityMap: toDomainActivityMap(a.Map.ActivityMap),
			Polyline:    a.Map.Polyline,
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

func toDomainActivityStreams(streams []stravapkg.ActivityStream) []domain.ActivityStream {
	out := make([]domain.ActivityStream, len(streams))
	for i, s := range streams {
		out[i] = domain.ActivityStream{
			Type:         s.Type,
			Data:         s.Data,
			SeriesType:   s.SeriesType,
			OriginalSize: s.OriginalSize,
			Resolution:   s.Resolution,
		}
	}
	return out
}

func toDomainCredentials(t stravapkg.TokenResponse) domain.StravaCredentials {
	return domain.StravaCredentials{
		AthleteID:    t.Athlete.ID,
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		ExpiresAt:    time.Unix(int64(t.ExpiresAt), 0),
	}
}

func toDomainRefreshCredentials(t stravapkg.RefreshTokenResponse) domain.StravaCredentials {
	return domain.StravaCredentials{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		ExpiresAt:    time.Unix(int64(t.ExpiresAt), 0),
	}
}
