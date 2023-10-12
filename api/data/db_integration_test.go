package data

import (
	"os"
	"testing"
	"time"
)

func TestDBIntegration(t *testing.T) {
	testDBURL := os.Getenv("DATABASE_URL")

	err := Migrate(testDBURL)
	if err != nil {
		t.Skip("skipping integration test")
	}

	db, err := NewDB(testDBURL)
	if err != nil {
		t.Skip("skipping integration test")
	}

	_, err = db.CreateUser("test", "test@test.com", "hashedpassword")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	_, err = db.GetUser("test")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	_, err = db.CreateArtist("Olivia Rodrigo")
	if err != nil {
		t.Error(err)
	}

	a, err := db.GetArtist("Olivia Rodrigo")
	if err != nil {
		t.Error(err)
	}

	release, _ := time.Parse("02/01/2006", "09/08/2023")
	p, err := db.CreateProject(CreateHash("GUTS", []string{"Olivia Rodrigo"}), "GUTS", []uint64{a.ID}, Album, release)
	if err != nil {
		t.Error(err)
	}

	_, err = db.GetProject(p.ID)
	if err != nil {
		t.Error(err)
	}

	track, err := db.CreateTrack(CreateHash("bad idea right?", []string{"Olivia Rodrigo"}), "bad idea right?", []uint64{a.ID})
	if err != nil {
		t.Error(err)
	}

	err = db.UpdateTrack(track.ID, p.ID, true)
	if err != nil {
		t.Error(err)
	}

	_, err = db.GetTrack(CreateHash("bad idea right?", []string{"Olivia Rodrigo"}))
	if err != nil {
		t.Error(err)
	}
}
