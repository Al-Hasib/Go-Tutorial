package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

// every file in this lesson folder connects to the same throwaway Postgres
// instance - see this folder's notes for how it was started.
const dsn = "postgres://postgres:postgres@localhost:55432/tutorial?sslmode=disable"

func main() {

	db, err := sql.Open("postgres", dsn)
	must(err)
	defer db.Close()
	must(db.Ping())

	// ---------- SETTING UP TABLES FOR THIS LESSON ----------

	//dropped and recreated every run, so this file can be re-run freely
	_, err = db.Exec(`DROP TABLE IF EXISTS sql_basics_books, sql_basics_authors`)
	must(err)
	_, err = db.Exec(`
		CREATE TABLE sql_basics_authors (
			id   SERIAL PRIMARY KEY,
			name TEXT NOT NULL
		)`)
	must(err)
	_, err = db.Exec(`
		CREATE TABLE sql_basics_books (
			id        SERIAL PRIMARY KEY,
			title     TEXT NOT NULL,
			author_id INTEGER REFERENCES sql_basics_authors(id),
			pages     INTEGER NOT NULL
		)`)
	must(err)

	// ---------- INSERT: ADDING ROWS ----------

	_, err = db.Exec(`INSERT INTO sql_basics_authors (name) VALUES ('Robert Martin'), ('Kathleen Booth')`)
	must(err)

	var martinID, boothID int
	must(db.QueryRow(`SELECT id FROM sql_basics_authors WHERE name = 'Robert Martin'`).Scan(&martinID))
	must(db.QueryRow(`SELECT id FROM sql_basics_authors WHERE name = 'Kathleen Booth'`).Scan(&boothID))

	_, err = db.Exec(`
		INSERT INTO sql_basics_books (title, author_id, pages) VALUES
			('Clean Code', $1, 464),
			('Clean Architecture', $1, 432),
			('Assembly Language Programming', $2, 200)`,
		martinID, boothID)
	must(err)

	// ---------- SELECT: READING ROWS ----------

	fmt.Println("all books:")
	printQuery(db, `SELECT id, title, pages FROM sql_basics_books ORDER BY id`)

	// ---------- WHERE: FILTERING ----------

	fmt.Println("\nbooks over 400 pages:")
	printQuery(db, `SELECT title, pages FROM sql_basics_books WHERE pages > 400`)

	// ---------- ORDER BY + LIMIT ----------

	fmt.Println("\nthe 2 shortest books:")
	printQuery(db, `SELECT title, pages FROM sql_basics_books ORDER BY pages ASC LIMIT 2`)

	// ---------- JOIN: COMBINING RELATED TABLES ----------

	fmt.Println("\nbooks with their author's name:")
	printQuery(db, `
		SELECT b.title, a.name
		FROM sql_basics_books b
		JOIN sql_basics_authors a ON b.author_id = a.id
		ORDER BY b.title`)

	// ---------- GROUP BY + AGGREGATE FUNCTIONS ----------

	fmt.Println("\nbook count and average pages per author:")
	printQuery(db, `
		SELECT a.name, COUNT(*), ROUND(AVG(b.pages))
		FROM sql_basics_books b
		JOIN sql_basics_authors a ON b.author_id = a.id
		GROUP BY a.name
		ORDER BY a.name`)

	// ---------- UPDATE: MODIFYING EXISTING ROWS ----------

	_, err = db.Exec(`UPDATE sql_basics_books SET pages = 470 WHERE title = 'Clean Code'`)
	must(err)
	fmt.Println("\nafter updating Clean Code's page count:")
	printQuery(db, `SELECT title, pages FROM sql_basics_books WHERE title = 'Clean Code'`)

	// ---------- DELETE: REMOVING ROWS ----------

	_, err = db.Exec(`DELETE FROM sql_basics_books WHERE title = 'Assembly Language Programming'`)
	must(err)
	fmt.Println("\nafter deleting one book:")
	printQuery(db, `SELECT title FROM sql_basics_books ORDER BY title`)
}

// printQuery runs a query and prints every row, generically - it doesn't
// know the column types ahead of time, unlike normal Go code, which is why
// it scans into `any` rather than typed variables. the dedicated
// database/sql lesson shows the normal, typed way to scan a row.
func printQuery(db *sql.DB, query string) {
	rows, err := db.Query(query)
	must(err)
	defer rows.Close()

	cols, err := rows.Columns()
	must(err)

	for rows.Next() {
		values := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range values {
			ptrs[i] = &values[i]
		}
		must(rows.Scan(ptrs...))

		parts := make([]string, len(cols))
		for i, v := range values {
			//NUMERIC/DECIMAL columns (like ROUND(AVG(...)) below) come back
			//from lib/pq as []byte, not a Go number - %v on a raw []byte
			//would print its byte values instead of the text they spell out
			if b, ok := v.([]byte); ok {
				parts[i] = string(b)
			} else {
				parts[i] = fmt.Sprintf("%v", v)
			}
		}
		fmt.Println(" ", strings.Join(parts, " | "))
	}
	must(rows.Err())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
