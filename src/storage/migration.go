package storage

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
)

//go:embed migrations/V01__initial.sql
var v01 string

//go:embed migrations/V02__rss_items.sql
var v02 string

//go:embed migrations/V03__drop_rss_items.sql
var v03 string

var migrations = []func(*sql.Tx) error{
	m01_initial,
	m02_rss_items,
	m03_drop_rss_items,
}
var maxVersion = int64(len(migrations))

func migrate(db *sql.DB) error {
	var version int64
	err := db.QueryRow("PRAGMA user_version").Scan(&version)
	if err != nil {
		return err
	}

	for v := version + 1; v <= maxVersion; v++ {
		log.Printf("[migration:%d] starting", v)
		if err = migrateVersion(v, db); err != nil {
			return err
		}
		log.Printf("[migration:%d] done", v)
	}

	return nil
}

func migrateVersion(v int64, db *sql.DB) error {
	var err error
	var tx *sql.Tx
	migratefunc := migrations[v-1]
	if tx, err = db.Begin(); err != nil {
		log.Printf("[migration:%d] failed to start transaction", v)
		return err
	}
	if err = migratefunc(tx); err != nil {
		log.Printf("[migration:%d] failed to migrate", v)
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec(fmt.Sprintf("pragma user_version = %d", v)); err != nil {
		log.Printf("[migration:%d] failed to bump version", v)
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		log.Printf("[migration:%d] failed to commit changes", v)
		return err
	}
	return nil
}

func m01_initial(tx *sql.Tx) error {
	_, err := tx.Exec(v01)
	return err
}

func m02_rss_items(tx *sql.Tx) error {
	_, err := tx.Exec(v02)
	return err
}

func m03_drop_rss_items(tx *sql.Tx) error {
	_, err := tx.Exec(v03)
	return err
}
