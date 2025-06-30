package migrations

import (
	"errors"

	"github.com/golang-migrate/migrate"
)

func RunPgMigrations(dsn string) error {
	if dsn == "" {
		return errors.New("no DSN provided")
	}

	path := ""
	// переделать на источник iofs
	// https://github.com/golang-migrate/migrate/blob/master/source/iofs/example_test.go
	m, err := migrate.New(
		path,
		dsn,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
