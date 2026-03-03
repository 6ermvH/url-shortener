package migrator

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/6ermvH/url-shortener/migrations"
)

func Run(db *sql.DB, version uint) error {
	sourceDriver, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("create migration source: %w", err)
	}

	dbDriver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		return fmt.Errorf("create migration driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "pgx5", dbDriver)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}

	if version != 0 {
		err = m.Migrate(version)
	} else {
		err = m.Up()
	}

	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
