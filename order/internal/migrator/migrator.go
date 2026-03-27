package migrator

import (
	"database/sql"

	"github.com/pressly/goose"
)

type Migrator struct {
	db            *sql.DB
	migrationsDir string
}

func NewMigrator(db *sql.DB, dir string) *Migrator {
	return &Migrator{
		db:            db,
		migrationsDir: dir,
	}
}

func (m *Migrator) Up() error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(m.db, m.migrationsDir)
}
