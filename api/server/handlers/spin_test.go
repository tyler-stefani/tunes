package handlers

import (
	"encoding/json"
	"testing"
	"time"

	"tunes-service/data"
)

type cacheMock struct {
	get func(string) string
	put func(string, string)
}

func (c *cacheMock) Get(key string) string {
	return c.get(key)
}

func (c *cacheMock) Put(key string, json string) {
	c.put(key, json)
}

type dbMock struct {
	getArtist     func(string) (data.Artist, error)
	createArtist  func(string) (data.Artist, error)
	getTrack      func(uint64) (data.Track, error)
	createTrack   func(uint64, string, []uint64) (data.Track, error)
	getProject    func(uint64) (data.Project, error)
	createProject func(uint64, string, []uint64, data.ProjectType, time.Time) (data.Project, error)
	createSpin    func(time.Time, uint64, uint64) (data.Spin, error)
	updateTrack   func(uint64, uint64, bool) error
}

func (d *dbMock) GetArtist(key string) (data.Artist, error) {
	return d.getArtist(key)
}

func (d *dbMock) CreateArtist(key string) (data.Artist, error) {
	return d.createArtist(key)
}

func (d *dbMock) GetTrack(key uint64) (data.Track, error) {
	return d.getTrack(key)
}

func (d *dbMock) CreateTrack(key uint64, title string, artistIDs []uint64) (data.Track, error) {
	return d.createTrack(key, title, artistIDs)
}

func (d *dbMock) GetProject(key uint64) (data.Project, error) {
	return d.getProject(key)
}

func (d *dbMock) CreateProject(key uint64, title string, artistIDs []uint64, projectType data.ProjectType, release time.Time) (data.Project, error) {
	return d.createProject(key, title, artistIDs, projectType, release)
}

func (d *dbMock) CreateSpin(time time.Time, userID uint64, trackID uint64) (data.Spin, error) {
	return d.createSpin(time, userID, trackID)
}

func (d *dbMock) UpdateTrack(trackID uint64, projectID uint64, isPrimary bool) error {
	return d.updateTrack(trackID, projectID, isPrimary)
}

func TestHandleSpin(t *testing.T) {
	release, _ := time.Parse("02/01/2006", "09/08/2023")
	spinTime := time.Now()
	tests := []struct {
		name     string
		input    SpinRequest
		db       *dbMock
		cache    *cacheMock
		expected data.Spin
	}{
		{
			"Should write new track and project",
			SpinRequest{
				1,
				spinTime,
				"bad idea right?",
				[]string{
					"Olivia Rodrigo",
				},
				"GUTS",
				[]string{
					"Olivia Rodrigo",
				},
				string(data.Album),
				release,
			},
			&dbMock{
				func(string) (data.Artist, error) {
					return data.Artist{}, nil
				},
				func(key string) (data.Artist, error) {
					return data.Artist{ID: 1, Name: key}, nil
				},
				func(uint64) (data.Track, error) {
					return data.Track{}, nil
				},
				func(key uint64, title string, artistIDs []uint64) (data.Track, error) {
					return data.Track{
						Title:            title,
						ProjectIDs:       []uint64{},
						PrimaryProjectID: 0,
					}, nil
				},
				func(uint64) (data.Project, error) {
					return data.Project{}, nil
				},
				func(key uint64, title string, artistIDs []uint64, projectType data.ProjectType, release time.Time) (data.Project, error) {
					return data.Project{
						Title:   title,
						Form:    projectType,
						Release: release,
					}, nil
				},
				func(time time.Time, userID uint64, trackID uint64) (data.Spin, error) {
					return data.Spin{
						ID:      uint(1),
						UserID:  uint(userID),
						Time:    time,
						TrackID: uint(trackID),
					}, nil
				},
				func(trackID uint64, projectID uint64, isPrimary bool) error {
					return nil
				},
			},
			&cacheMock{
				func(string) string {
					return ""
				},
				func(string, string) {
				},
			},
			data.Spin{
				ID:      1,
				UserID:  1,
				Time:    spinTime,
				TrackID: 908849726797084829,
			},
		},
		{
			"Should not hit database when cached",
			SpinRequest{
				1,
				spinTime,
				"bad idea right?",
				[]string{
					"Olivia Rodrigo",
				},
				"GUTS",
				[]string{
					"Olivia Rodrigo",
				},
				string(data.Album),
				release,
			},
			&dbMock{
				func(string) (data.Artist, error) {
					t.FailNow()
					return data.Artist{}, nil
				},
				func(key string) (data.Artist, error) {
					t.FailNow()
					return data.Artist{}, nil
				},
				func(uint64) (data.Track, error) {
					t.FailNow()
					return data.Track{}, nil
				},
				func(uint64, string, []uint64) (data.Track, error) {
					t.FailNow()
					return data.Track{}, nil
				},
				func(uint64) (data.Project, error) {
					t.FailNow()
					return data.Project{}, nil
				},
				func(uint64, string, []uint64, data.ProjectType, time.Time) (data.Project, error) {
					t.FailNow()
					return data.Project{}, nil
				},
				func(time time.Time, userID uint64, trackID uint64) (data.Spin, error) {
					return data.Spin{
						ID:      uint(1),
						UserID:  uint(userID),
						Time:    time,
						TrackID: uint(trackID),
					}, nil
				},
				func(trackID uint64, projectID uint64, isPrimary bool) error {
					return nil
				},
			},
			&cacheMock{
				func(key string) string {
					if key == "a-Olivia Rodrigo" {
						j, _ := json.Marshal(data.Artist{
							ID:   1,
							Name: "Olivia Rodrigo",
						})
						return string(j)
					} else if key == "t-908849726797084829" {
						j, _ := json.Marshal(data.Track{
							Title:            "bad idea right?",
							ProjectIDs:       []uint64{2371983535859309014},
							PrimaryProjectID: 2371983535859309014,
						})
						return string(j)
					} else if key == "p-2371983535859309014" {
						j, _ := json.Marshal(data.Project{
							Title:   "GUTS",
							Form:    data.Album,
							Release: release,
						})
						return string(j)
					}
					return ""
				},
				func(string, string) {
				},
			},
			data.Spin{
				ID:      1,
				UserID:  1,
				Time:    spinTime,
				TrackID: 908849726797084829,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := HandleSpin(tt.input, tt.db, tt.cache)
			if actual != tt.expected {
				t.Fatalf("expected %+v but got %+v", tt.expected, actual)
			}
		})
	}
}
