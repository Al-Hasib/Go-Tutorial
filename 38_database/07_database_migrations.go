package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const dsn = "postgres://postgres:postgres@localhost:55432/tutorial?sslmode=disable"

// migration is intentionally tiny - real tools (golang-migrate, goose,
// atlas, ...) do this with separate .sql files, checksums, and a matching
// "down" migration for reverting. this hand-rolled version exists purely
// to make the underlying idea clear before you reach for one of those.
type migration struct {
	version     int
	description string
	up          string
}

var migrations = []migration{
	{1, "create users table", `
		CREATE TABLE migrate_users (
			id   SERIAL PRIMARY KEY,
			name TEXT NOT NULL
		)`},
	{2, "add email column to users", `
		ALTER TABLE migrate_users ADD COLUMN email TEXT`},
	{3, "create posts table", `
		CREATE TABLE migrate_posts (
			id      SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES migrate_users(id),
			title   TEXT NOT NULL
		)`},
}

func main() {

	db, err := sql.Open("postgres", dsn)
	must(err)
	defer db.Close()
	must(db.Ping())

	//start from a clean slate so this lesson can be re-run repeatedly
	_, err = db.Exec(`DROP TABLE IF EXISTS migrate_posts, migrate_users, schema_migrations`)
	must(err)

	fmt.Println("--- first run ---")
	runMigrations(db)

	fmt.Println("\n--- running again (should apply nothing new) ---")
	runMigrations(db)
}

// runMigrations applies every migration whose version isn't already
// recorded in schema_migrations, in order, each inside its own transaction -
// if one fails partway through, ITS changes roll back while every
// migration before it stays applied and recorded.
func runMigrations(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    INTEGER PRIMARY KEY,
			applied_at TIMESTAMPTZ DEFAULT now()
		)`)
	must(err)

	applied := map[int]bool{}
	rows, err := db.Query(`SELECT version FROM schema_migrations`)
	must(err)
	for rows.Next() {
		var v int
		must(rows.Scan(&v))
		applied[v] = true
	}
	rows.Close()

	for _, m := range migrations {
		if applied[m.version] {
			fmt.Printf("  v%d (%s) - already applied, skipping\n", m.version, m.description)
			continue
		}

		tx, err := db.Begin()
		must(err)

		if _, err := tx.Exec(m.up); err != nil {
			tx.Rollback()
			panic(fmt.Errorf("migration v%d failed: %w", m.version, err))
		}
		if _, err := tx.Exec(`INSERT INTO schema_migrations (version) VALUES ($1)`, m.version); err != nil {
			tx.Rollback()
			panic(err)
		}
		must(tx.Commit())
		fmt.Printf("  v%d (%s) - applied\n", m.version, m.description)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
