package main

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

const dsn = "postgres://postgres:postgres@localhost:55432/tutorial?sslmode=disable"

func main() {

	db, err := sql.Open("postgres", dsn)
	must(err)
	defer db.Close()
	must(db.Ping())

	// ---------- POSTGRES VERSION + SERVER INFO ----------

	var version string
	must(db.QueryRow(`SELECT version()`).Scan(&version))
	fmt.Println("connected to:", version)

	// ---------- LISTING TABLES VIA information_schema ----------

	//information_schema is a standard SQL feature Postgres implements -
	//useful for introspecting a database without a special client tool
	rows, err := db.Query(`
		SELECT table_name FROM information_schema.tables
		WHERE table_schema = 'public' ORDER BY table_name`)
	must(err)
	fmt.Println("\ntables in this database:")
	for rows.Next() {
		var name string
		must(rows.Scan(&name))
		fmt.Println(" -", name)
	}
	rows.Close()

	_, err = db.Exec(`DROP TABLE IF EXISTS pg_features_products`)
	must(err)
	_, err = db.Exec(`
		CREATE TABLE pg_features_products (
			id         SERIAL PRIMARY KEY,
			sku        TEXT UNIQUE NOT NULL,
			name       TEXT NOT NULL,
			tags       TEXT[],
			metadata   JSONB,
			created_at TIMESTAMPTZ DEFAULT now()
		)`)
	must(err)

	// ---------- RETURNING: GET DATA BACK FROM AN INSERT/UPDATE/DELETE ----------

	//RETURNING is a Postgres extension (a few other databases have it too) -
	//saves a second round-trip just to fetch what you inserted, like an
	//auto-generated id
	var newID int
	must(db.QueryRow(`
		INSERT INTO pg_features_products (sku, name, tags, metadata)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		"SKU-1", "Widget", pq.Array([]string{"tools", "hardware"}), `{"color":"blue","weight_kg":1.2}`,
	).Scan(&newID))
	fmt.Println("\ninserted product with id:", newID)

	// ---------- ARRAYS ----------

	//pq.Array bridges a Go slice to a Postgres array column, both directions
	var tags []string
	must(db.QueryRow(`SELECT tags FROM pg_features_products WHERE id = $1`, newID).Scan(pq.Array(&tags)))
	fmt.Println("tags:", tags)

	// ---------- JSONB ----------

	//JSONB stores JSON in a queryable, indexable binary form - -> extracts a
	//field as JSON, ->> extracts it as text
	var color string
	must(db.QueryRow(`SELECT metadata->>'color' FROM pg_features_products WHERE id = $1`, newID).Scan(&color))
	fmt.Println("metadata color:", color)

	// ---------- ON CONFLICT: UPSERT IN ONE STATEMENT ----------

	//without ON CONFLICT this needs SELECT-then-INSERT-or-UPDATE (2 round
	//trips, and a race between them) - Postgres does it atomically in one
	_, err = db.Exec(`
		INSERT INTO pg_features_products (sku, name)
		VALUES ($1, $2)
		ON CONFLICT (sku) DO UPDATE SET name = EXCLUDED.name`,
		"SKU-1", "Widget (updated via upsert)")
	must(err)
	var updatedName string
	must(db.QueryRow(`SELECT name FROM pg_features_products WHERE sku = $1`, "SKU-1").Scan(&updatedName))
	fmt.Println("after upsert:", updatedName)

	// ---------- POSTGRES ERROR CODES ----------

	//violating the UNIQUE constraint on sku on purpose, then inspecting the
	//specific error code Postgres sent back
	_, err = db.Exec(`INSERT INTO pg_features_products (sku, name) VALUES ($1, $2)`, "SKU-1", "Duplicate")
	if pqErr, ok := err.(*pq.Error); ok {
		fmt.Println("\ngot a Postgres error - code:", pqErr.Code, "name:", pqErr.Code.Name())
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
