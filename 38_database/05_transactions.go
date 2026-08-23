package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

const dsn = "postgres://postgres:postgres@localhost:55432/tutorial?sslmode=disable"

func main() {

	db, err := sql.Open("postgres", dsn)
	must(err)
	defer db.Close()
	must(db.Ping())

	_, err = db.Exec(`DROP TABLE IF EXISTS tx_accounts`)
	must(err)
	_, err = db.Exec(`
		CREATE TABLE tx_accounts (
			id      SERIAL PRIMARY KEY,
			name    TEXT NOT NULL,
			balance INTEGER NOT NULL
		)`)
	must(err)
	_, err = db.Exec(`INSERT INTO tx_accounts (name, balance) VALUES ('Alice', 100), ('Bob', 20)`)
	must(err)

	// ---------- A SUCCESSFUL TRANSFER ----------

	must(transfer(db, "Alice", "Bob", 30))
	fmt.Println("after a normal transfer:")
	printBalances(db)

	// ---------- WHY A TRANSACTION MATTERS: ATOMICITY ----------

	//this transfer moves 1000 from Bob (who only has 50) to Alice - the
	//debit would succeed as a standalone statement, but the business-rule
	//check below fails it INSIDE the transaction, so both the debit and the
	//credit are rolled back together. without a transaction wrapping both
	//statements, the debit could already have committed before the check
	//even ran, leaving Bob's balance wrong with no matching credit anywhere.
	err = transfer(db, "Bob", "Alice", 1000)
	fmt.Println("\nattempted an overdraft transfer, got error:", err)
	fmt.Println("balances (unchanged - the failed transfer left no partial trace):")
	printBalances(db)

	fmt.Println("\ndone")
}

// transfer moves amount from `from` to `to`, atomically: either both the
// debit and the credit happen, or neither does.
func transfer(db *sql.DB, from, to string, amount int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	//if this function returns before Commit, tx.Rollback() undoes
	//everything done in this transaction so far - a safe default, since
	//calling Rollback after a successful Commit is just a harmless no-op
	defer tx.Rollback()

	var balance int
	if err := tx.QueryRow(`SELECT balance FROM tx_accounts WHERE name = $1`, from).Scan(&balance); err != nil {
		return err
	}
	if balance < amount {
		return errors.New("insufficient funds") // returning here triggers the deferred Rollback
	}

	if _, err := tx.Exec(`UPDATE tx_accounts SET balance = balance - $1 WHERE name = $2`, amount, from); err != nil {
		return err
	}
	if _, err := tx.Exec(`UPDATE tx_accounts SET balance = balance + $1 WHERE name = $2`, amount, to); err != nil {
		return err
	}

	return tx.Commit() // only now do both updates become permanent
}

func printBalances(db *sql.DB) {
	rows, err := db.Query(`SELECT name, balance FROM tx_accounts ORDER BY name`)
	must(err)
	defer rows.Close()
	for rows.Next() {
		var name string
		var balance int
		must(rows.Scan(&name, &balance))
		fmt.Printf("  %s: %d\n", name, balance)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
