package data

import "time"

type DB interface {
	UserDB
	TunesDB
}

type UserDB interface {
	GetUser(nameOrEmail string) (User, error)
	CreateUser(name, email, password string) (User, error)
}

type TunesDB interface {
	GetArtist(name string) (Artist, error)
	CreateArtist(name string) (Artist, error)
	GetTrack(key uint64) (Track, error)
	CreateTrack(key uint64, title string, artistIDs []uint64) (Track, error)
	GetProject(key uint64) (Project, error)
	CreateProject(key uint64, title string, artistIDs []uint64, form ProjectType, release time.Time) (Project, error)
	CreateSpin(t time.Time, userID uint64, trackID uint64) (Spin, error)
	UpdateTrack(key uint64, projectID uint64, isPrimary bool) error
}

type AuthDB interface {
	WriteRefreshToken(id string, expires time.Time) (bool, error)
	FindRefreshToken(id string) (RefreshToken, error)
}
