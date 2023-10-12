package data

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type PGDB struct {
	db *pgx.Conn
}

func NewDB(url string) (*PGDB, error) {
	db, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return &PGDB{}, err
	} else {
		return &PGDB{
			db,
		}, nil
	}
}

func (pg *PGDB) CreateUser(name string, email string, password string) (User, error) {
	const stmt = `INSERT INTO "user" (name, email, password) VALUES ($1, $2, $3) RETURNING (id, name, email, password)`

	row := pg.db.QueryRow(context.Background(), stmt, name, email, password)

	var u User
	if err := row.Scan(&u); err != nil {
		return User{}, err
	}

	return u, nil
}

func (pg *PGDB) CreateArtist(name string) (Artist, error) {
	const stmt = `INSERT INTO artist (name) VALUES ($1) RETURNING (id, name)`

	row := pg.db.QueryRow(context.Background(), stmt, name)

	var a Artist
	if err := row.Scan(&a); err != nil {
		return Artist{}, fmt.Errorf("error inserting artist: %w", err)
	}

	return a, nil
}

func (pg *PGDB) CreateProject(key uint64, title string, artistIDs []uint64, form ProjectType, release time.Time) (Project, error) {
	const stmt = `INSERT INTO project (id, title, form, release) VALUES ($1, $2, $3, $4) RETURNING (id, title, form, release)`
	const junctionInsert = `INSERT INTO artist_project (artist_id, project_id) VALUES ($1, $2)`

	row := pg.db.QueryRow(context.Background(), stmt, key, title, form, release)

	var p Project
	if err := row.Scan(&p); err != nil {
		return Project{}, fmt.Errorf("error inserting project: %w", err)
	}

	for _, artistID := range artistIDs {
		pg.db.Exec(context.Background(), junctionInsert, artistID, key)
	}

	return p, nil
}

func (pg *PGDB) CreateSpin(t time.Time, userID uint64, trackID uint64) (Spin, error) {
	const stmt = `INSERT INTO spin (time, user_id, track_id) VALUES ($1, $2, $3) RETURNING (time, user_id, track_id)`

	row := pg.db.QueryRow(context.Background(), stmt, t, userID, trackID)

	var s Spin
	if err := row.Scan(&s); err != nil {
		return Spin{}, fmt.Errorf("error inserting spin: %w", err)
	}

	return s, nil
}

func (pg *PGDB) CreateTrack(key uint64, title string, artistIDs []uint64) (Track, error) {
	const stmt = `INSERT INTO track (id, title) VALUES ($1, $2) RETURNING (id, title)`
	const junctionInsert = `INSERT INTO artist_track (artist_id, track_id) VALUES ($1, $2)`

	row := pg.db.QueryRow(context.Background(), stmt, key, title)

	var t Track
	if err := row.Scan(&t); err != nil {
		return Track{}, fmt.Errorf("error inserting track: %w", err)
	}

	for _, artistID := range artistIDs {
		pg.db.Exec(context.Background(), junctionInsert, artistID, key)
	}

	return t, nil
}

func (pg *PGDB) GetUser(nameOrEmail string) (User, error) {
	const stmt = `SELECT (id, name, email, password) FROM "user" WHERE name=$1 OR email=$1`

	row := pg.db.QueryRow(context.Background(), stmt, nameOrEmail)

	var u User
	if err := row.Scan(&u); err != nil {
		return User{}, fmt.Errorf("error selecting user: %w", err)
	}

	return u, nil
}

func (pg *PGDB) GetArtist(name string) (Artist, error) {
	const stmt = `SELECT (id, name) FROM artist WHERE name=$1`

	row := pg.db.QueryRow(context.Background(), stmt, name)

	var a Artist
	if err := row.Scan(&a); err != nil {
		return Artist{}, fmt.Errorf("error selecting artist: %w", err)
	}

	return a, nil
}

func (pg *PGDB) GetProject(key uint64) (Project, error) {
	const stmt = `SELECT (id, title, form, release) FROM project WHERE id=$1`

	row := pg.db.QueryRow(context.Background(), stmt, key)

	var p Project
	if err := row.Scan(&p); err != nil {
		return Project{}, fmt.Errorf("error selecting project: %w", err)
	}

	return p, nil
}

// TODO: maybe make this query work for tracks with no projects
func (pg *PGDB) GetTrack(key uint64) (Track, error) {
	stmt := `SELECT t.id, t.title, t.primary_project_id, array_agg(p.id) AS project_ids
	FROM track t
	JOIN project_track pt ON t.id = pt.track_id
	JOIN project p ON pt.project_id = p.id
	WHERE t.id = $1
	GROUP BY t.id, t.title, t.primary_project_id;`

	row := pg.db.QueryRow(context.Background(), stmt, key)

	var t Track
	if err := row.Scan(&t.ID, &t.Title, &t.PrimaryProjectID, &t.ProjectIDs); err != nil {
		return Track{}, fmt.Errorf("error selecting track: %w", err)
	}

	return t, nil
}

func (pg *PGDB) UpdateTrack(key uint64, projectID uint64, isPrimary bool) error {
	const junctionInsert = `INSERT INTO project_track (project_id, track_id) VALUES ($2, $1)`
	const primaryProjectUpdate = `UPDATE track SET primary_project_id=$2 WHERE id=$1`

	pg.db.Exec(context.Background(), junctionInsert, key, projectID)

	if isPrimary {
		pg.db.Exec(context.Background(), primaryProjectUpdate, key, projectID)
	}

	return nil
}
