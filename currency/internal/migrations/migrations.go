package migrations

import (
	"embed"
	"errors"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var migrationFiles embed.FS

func RunPgMigrations(dsn string) error {
	if dsn == "" {
		return errors.New("no DSN provided")
	}

	sub, err := fs.Sub(migrationFiles, ".")
	if err != nil {
		return err
	}
	d, err := iofs.New(sub, ".")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
