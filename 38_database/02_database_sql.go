package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // registers the "postgres" driver with database/sql - imported only for its side effect
)

const dsn = "postgres://postgres:postgres@localhost:55432/tutorial?sslmode=disable"

type Person struct {
	ID       int
	Name     string
	Nickname sql.NullString // a NULL-able column needs a Null* wrapper, not a plain string
}

func main() {

	// ---------- sql.Open IS LAZY ----------

	//sql.Open just validates the DSN's shape and prepares a *sql.DB handle -
	//it does NOT connect to anything yet. Ping (or the first real query) is
	//what actually opens a connection and can surface "wrong host/password".
	db, err := sql.Open("postgres", dsn)
	must(err)
	defer db.Close()

	must(db.Ping())

	_, err = db.Exec(`DROP TABLE IF EXISTS dbsql_people`)
	must(err)
	_, err = db.Exec(`
		CREATE TABLE dbsql_people (
			id       SERIAL PRIMARY KEY,
			name     TEXT NOT NULL,
			nickname TEXT
		)`)
	must(err)

	// ---------- PARAMETERIZED QUERIES: NEVER STRING-CONCATENATE USER INPUT ----------

	//WRONG (never do this - shown only as a comment, never actually run):
	//   query := "INSERT INTO dbsql_people (name) VALUES ('" + userInput + "')"
	//   a name like "x'); DROP TABLE dbsql_people; --" would execute AS SQL,
	//   not as harmless data. that's SQL injection.
	//RIGHT: let the driver substitute values safely via $1, $2, ... placeholders
	_, err = db.Exec(`INSERT INTO dbsql_people (name, nickname) VALUES ($1, $2)`, "Alice", "Ali")
	must(err)
	_, err = db.Exec(`INSERT INTO dbsql_people (name, nickname) VALUES ($1, $2)`, "Bob", nil) // no nickname
	must(err)

	// ---------- PREPARED STATEMENTS: PARSE ONCE, RUN MANY TIMES ----------

	//useful when running the same shaped query repeatedly - the database
	//parses/plans it once instead of redoing that work on every call
	stmt, err := db.Prepare(`INSERT INTO dbsql_people (name, nickname) VALUES ($1, $2)`)
	must(err)
	defer stmt.Close()
	for _, n := range []string{"Carol", "Dave"} {
		_, err = stmt.Exec(n, nil)
		must(err)
	}

	// ---------- Exec VS Query VS QueryRow ----------

	//Exec: for statements that don't return rows (INSERT/UPDATE/DELETE/DDL) -
	//gives back a Result you can ask for rows affected
	res, err := db.Exec(`UPDATE dbsql_people SET nickname = $1 WHERE name = $2`, "Bobby", "Bob")
	must(err)
	affected, _ := res.RowsAffected()
	fmt.Println("rows affected by UPDATE:", affected)

	//QueryRow: for exactly one row expected - Scan reports sql.ErrNoRows if there wasn't one
	var count int
	must(db.QueryRow(`SELECT COUNT(*) FROM dbsql_people`).Scan(&count))
	fmt.Println("total people:", count)

	//Query: for zero or more rows - always range with rows.Next(), always
	//defer rows.Close(), and check rows.Err() after the loop (a failure
	//mid-stream shows up there, not as a panic and not from Next() itself)
	rows, err := db.Query(`SELECT id, name, nickname FROM dbsql_people ORDER BY id`)
	must(err)
	defer rows.Close()

	fmt.Println("\npeople:")
	for rows.Next() {
		var p Person
		//nickname can be SQL NULL - scanning straight into a string would
		//fail on a NULL row; sql.NullString has a .Valid flag for this
		must(rows.Scan(&p.ID, &p.Name, &p.Nickname))
		if p.Nickname.Valid {
			fmt.Printf("  #%d %s (aka %s)\n", p.ID, p.Name, p.Nickname.String)
		} else {
			fmt.Printf("  #%d %s (no nickname)\n", p.ID, p.Name)
		}
	}
	must(rows.Err()) // always check this - Next() returning false doesn't mean "no error happened"

	// ---------- CONTEXT-AWARE QUERIES: BOUNDING HOW LONG A QUERY CAN RUN ----------

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var name string
	must(db.QueryRowContext(ctx, `SELECT name FROM dbsql_people WHERE id = $1`, 1).Scan(&name))
	fmt.Println("\nQueryRowContext result:", name)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
