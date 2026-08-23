package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

const dsn = "postgres://postgres:postgres@localhost:55432/tutorial?sslmode=disable"

type Task struct {
	ID    int
	Title string
	Done  bool
}

// TaskRepo bundles all the CRUD operations for one table behind a small,
// testable API - callers work with Task values, never write SQL themselves.
type TaskRepo struct {
	db *sql.DB
}

// ---------- CREATE ----------
func (r *TaskRepo) Create(title string) (Task, error) {
	var t Task
	t.Title = title
	err := r.db.QueryRow(
		`INSERT INTO crud_tasks (title, done) VALUES ($1, false) RETURNING id, done`,
		title,
	).Scan(&t.ID, &t.Done)
	return t, err
}

// ---------- READ (one) ----------
func (r *TaskRepo) Get(id int) (Task, error) {
	var t Task
	err := r.db.QueryRow(`SELECT id, title, done FROM crud_tasks WHERE id = $1`, id).
		Scan(&t.ID, &t.Title, &t.Done)
	if errors.Is(err, sql.ErrNoRows) {
		return Task{}, fmt.Errorf("task %d not found", id)
	}
	return t, err
}

// ---------- READ (many) ----------
func (r *TaskRepo) List() ([]Task, error) {
	rows, err := r.db.Query(`SELECT id, title, done FROM crud_tasks ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

// ---------- UPDATE ----------
func (r *TaskRepo) Update(id int, title string, done bool) error {
	res, err := r.db.Exec(`UPDATE crud_tasks SET title = $1, done = $2 WHERE id = $3`, title, done, id)
	if err != nil {
		return err
	}
	//an UPDATE with no matching row is NOT an error as far as SQL is
	//concerned - checking RowsAffected is how you turn "updated nothing"
	//into a real "not found" for the caller
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("task %d not found", id)
	}
	return nil
}

// ---------- DELETE ----------
func (r *TaskRepo) Delete(id int) error {
	res, err := r.db.Exec(`DELETE FROM crud_tasks WHERE id = $1`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("task %d not found", id)
	}
	return nil
}

func main() {

	db, err := sql.Open("postgres", dsn)
	must(err)
	defer db.Close()
	must(db.Ping())

	_, err = db.Exec(`DROP TABLE IF EXISTS crud_tasks`)
	must(err)
	_, err = db.Exec(`
		CREATE TABLE crud_tasks (
			id    SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			done  BOOLEAN NOT NULL DEFAULT false
		)`)
	must(err)

	repo := &TaskRepo{db: db}

	// CREATE
	t1, err := repo.Create("write the CRUD lesson")
	must(err)
	t2, err := repo.Create("review it")
	must(err)
	fmt.Println("created:", t1, t2)

	// READ (one)
	got, err := repo.Get(t1.ID)
	must(err)
	fmt.Println("\nfetched:", got)

	// READ (many)
	all, err := repo.List()
	must(err)
	fmt.Println("\nall tasks:", all)

	// UPDATE
	must(repo.Update(t1.ID, "write the CRUD lesson", true))
	got, err = repo.Get(t1.ID)
	must(err)
	fmt.Println("\nafter update:", got)

	// UPDATE on a missing row - a clean error, not a crash
	if err := repo.Update(9999, "x", true); err != nil {
		fmt.Println("\nupdating a missing task:", err)
	}

	// DELETE
	must(repo.Delete(t2.ID))
	all, err = repo.List()
	must(err)
	fmt.Println("\nafter deleting one task:", all)

	// READ a deleted row - also a clean error
	if _, err := repo.Get(t2.ID); err != nil {
		fmt.Println("\nfetching a deleted task:", err)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
