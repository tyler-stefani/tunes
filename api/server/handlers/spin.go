package handlers

import (
	"encoding/json"
	"slices"
	"strconv"
	"time"

	c "tunes-service/cache"
	d "tunes-service/data"
)

type SpinRequest struct {
	UserID             uint
	Time               time.Time
	TrackTitle         string
	TrackArtistNames   []string
	ProjectTitle       string
	ProjectArtistNames []string
	ProjectType        string
	ProjectRelese      time.Time
}

func HandleSpin(req SpinRequest, db d.TunesDB, cache c.Cache) d.Spin {
	trackHash := d.CreateHash(req.TrackTitle, req.TrackArtistNames)
	t := getTrack(trackHash, db, cache)
	if t.IsEmpty() {
		artistIDs := []uint64{}
		for _, artistName := range req.TrackArtistNames {
			a := getArtist(artistName, db, cache)
			if a.IsEmpty() {
				a, _ = db.CreateArtist(artistName)
			}
			artistIDs = append(artistIDs, a.ID)
		}
		t, _ = db.CreateTrack(trackHash, req.TrackTitle, artistIDs)
	}

	projectHash := d.CreateHash(req.ProjectTitle, req.ProjectArtistNames)
	p := getProject(projectHash, db, cache)
	if p.IsEmpty() {
		artistIDs := []uint64{}
		for _, artistName := range req.TrackArtistNames {
			a := getArtist(artistName, db, cache)
			if a.IsEmpty() {
				a, _ = db.CreateArtist(artistName)
			}
			artistIDs = append(artistIDs, a.ID)
		}
		p, _ = db.CreateProject(projectHash, req.ProjectTitle, artistIDs, d.ProjectType(req.ProjectType), req.ProjectRelese)
	}

	if !slices.Contains(t.ProjectIDs, projectHash) {
		primaryProject, _ := db.GetProject(t.PrimaryProjectID)

		db.UpdateTrack(trackHash, projectHash, primaryProject.IsLessPrimaryThan(&p))
	}

	s, _ := db.CreateSpin(req.Time, uint64(req.UserID), trackHash)
	return s
}

func getArtist(key string, db d.TunesDB, cache c.Cache) (a d.Artist) {
	cachedJSON := cache.Get("a-" + key)
	if cachedJSON != "" {
		json.Unmarshal([]byte(cachedJSON), &a)
		return
	}

	a, _ = db.GetArtist(key)
	if !a.IsEmpty() {
		j, _ := json.Marshal(a)
		cache.Put("a-"+key, string(j))
	}
	return
}

func getTrack(key uint64, db d.TunesDB, cache c.Cache) (t d.Track) {
	cachedJSON := cache.Get("t-" + strconv.FormatUint(key, 10))
	if cachedJSON != "" {
		json.Unmarshal([]byte(cachedJSON), &t)
		return
	}

	t, _ = db.GetTrack(key)
	if !t.IsEmpty() {
		j, _ := json.Marshal(t)
		cache.Put("t-"+strconv.FormatUint(key, 10), string(j))
	}
	return
}

func getProject(key uint64, db d.TunesDB, cache c.Cache) (p d.Project) {
	cachedJSON := cache.Get("p-" + strconv.FormatUint(key, 10))
	if cachedJSON != "" {
		json.Unmarshal([]byte(cachedJSON), &p)
		return
	}

	p, _ = db.GetProject(key)
	if !p.IsEmpty() {
		j, _ := json.Marshal(p)
		cache.Put("p-"+strconv.FormatUint(key, 10), string(j))
	}
	return
}
