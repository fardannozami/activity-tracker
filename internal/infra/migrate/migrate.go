package migrate

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigration(databaseUrl string) error {
	m, err := migrate.New("file://migrations", databaseUrl)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("database already up to date")
			return nil
		}
		return err
	}

	log.Println("database migration applied")
	return nil
}
