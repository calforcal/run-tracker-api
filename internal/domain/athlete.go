package domain

import "time"

// Bike is a piece of Strava gear.
type Bike struct {
	ID            string
	Name          string
	Primary       bool
	ResourceState int
	Distance      float64
}

// Shoe is a piece of Strava gear.
type Shoe struct {
	ID            string
	Name          string
	Primary       bool
	ResourceState int
	Distance      float64
}

// Athlete is a Strava athlete profile.
type Athlete struct {
	ID                    int64
	Username              string
	ResourceState         int
	Firstname             string
	Lastname              string
	City                  string
	State                 string
	Country               string
	Sex                   string
	Premium               bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
	BadgeTypeID           int
	ProfileMedium         string
	Profile               string
	Friend                *bool
	Follower              *bool
	FollowerCount         int
	FriendCount           int
	MutualFriendCount     int
	AthleteType           int
	DatePreference        string
	MeasurementPreference string
	FTP                   *float64
	Weight                float64
	Bikes                 []Bike
	Shoes                 []Shoe
}
