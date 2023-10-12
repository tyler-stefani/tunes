package data

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(url string) error {
	m, err := migrate.New("file:./migrations", url)
	if err != nil {
		return err
	}
	m.Up()
	return nil
}
