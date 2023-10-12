package data

import (
	"sort"
	"strings"
	"time"

	"github.com/cespare/xxhash/v2"
)

type User struct {
	ID       uint64
	Name     string
	Email    string
	Password string
}

type Artist struct {
	ID   uint64
	Name string
}

func (a *Artist) IsEmpty() bool {
	return a.ID == 0 && a.Name == ""
}

type ProjectType string

const (
	Album       ProjectType = "album"
	EP          ProjectType = "ep"
	Single      ProjectType = "single"
	Compilation ProjectType = "compilation"
)

type Project struct {
	ID      uint64
	Title   string
	Form    ProjectType
	Release time.Time
}

func (p *Project) IsEmpty() bool {
	return p.ID == 0 && p.Title == "" && p.Form == "" && p.Release == time.Time{}
}

func (a *Project) IsLessPrimaryThan(b *Project) bool {
	if a.Form == b.Form {
		return a.Release.Before(b.Release)
	}

	if a.Form == Album {
		return false
	} else if a.Form == Compilation {
		return true
	} else if a.Form == EP {
		if b.Form == Album {
			return true
		} else {
			return false
		}
	} else {
		if b.Form == Compilation {
			return false
		} else {
			return true
		}
	}
}

type Track struct {
	ID               uint64
	Title            string
	ProjectIDs       []uint64
	PrimaryProjectID uint64
}

func (t *Track) IsEmpty() bool {
	return t.ID == 0 && t.Title == "" && len(t.ProjectIDs) == 0 && t.PrimaryProjectID == 0
}

type Spin struct {
	ID      uint
	UserID  uint
	Time    time.Time
	TrackID uint
}

func CreateHash(title string, artistNames []string) uint64 {
	sort.Slice(artistNames, func(i, j int) bool {
		return artistNames[i] < artistNames[j]
	})
	prehash := strings.Join(append([]string{title}, artistNames...), "|")
	return xxhash.Sum64String(prehash)
}

type RefreshToken struct {
	ID         string `bson:"_id"`
	Expiration time.Time
}
