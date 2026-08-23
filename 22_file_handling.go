package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func main() {

	// everything below happens inside a temp folder so this lesson
	// doesn't leave files behind in your project folder
	dir, err := os.MkdirTemp("", "go-lesson-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir) // cleanup at the very end

	// ---------- FILE HANDLING ----------

	//building a file path the SAFE way - filepath.Join
	//(handles / vs \ for you, works on Windows and Linux/Mac)
	filePath := filepath.Join(dir, "notes.txt")
	fmt.Println("file path:", filePath)

	//creating a file with os.Create - makes an empty file and gives
	//back a *os.File you can write to (an existing file gets truncated)
	created, err := os.Create(filePath)
	if err != nil {
		fmt.Println("create error:", err)
	}
	created.WriteString("Hello, Go!\n")
	created.Close() // always close what you open, or use defer right after

	//writing a file in one call - simplest way when you just have the
	//full content ready (creates the file if needed, or overwrites it)
	err = os.WriteFile(filePath, []byte("Hello, Go!\n"), 0644)
	if err != nil {
		fmt.Println("write error:", err)
	}

	//reading a whole file in one call
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("read error:", err)
	}
	fmt.Println("read content:", string(data))

	//appending more text to an existing file
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("open error:", err)
	}
	f.WriteString("Second line\n")
	f.Close() // must close, or use defer f.Close() right after opening

	//reading line by line with bufio.Scanner - useful for big files
	f2, _ := os.Open(filePath)
	defer f2.Close()
	scanner := bufio.NewScanner(f2)
	fmt.Println("line by line:")
	for scanner.Scan() {
		fmt.Println(" ", scanner.Text())
	}

	//checking if a file exists
	info, err := os.Stat(filePath)
	if err != nil {
		fmt.Println("file does not exist")
	} else {
		fmt.Println("exists:", info.Name(), "size:", info.Size(), "bytes")
	}

	//checking a path that does NOT exist
	_, err = os.Stat(filepath.Join(dir, "missing.txt"))
	if os.IsNotExist(err) {
		fmt.Println("missing.txt really is missing")
	}

	//renaming (moving) a file
	renamedPath := filepath.Join(dir, "renamed.txt")
	os.Rename(filePath, renamedPath)
	fmt.Println("renamed to:", renamedPath)

	//deleting a file
	os.Remove(renamedPath)
	_, err = os.Stat(renamedPath)
	fmt.Println("deleted, exists now?", !os.IsNotExist(err))

	// ---------- DIRECTORIES & PATHS ----------

	//creating one directory
	subDir := filepath.Join(dir, "sub")
	os.Mkdir(subDir, 0755)

	//creating nested directories at once - Mkdir alone would fail here
	nestedDir := filepath.Join(dir, "a", "b", "c")
	err = os.MkdirAll(nestedDir, 0755)
	fmt.Println("created nested dirs:", err == nil)

	//putting a couple of files in place to list afterwards
	os.WriteFile(filepath.Join(dir, "one.txt"), []byte("1"), 0644)
	os.WriteFile(filepath.Join(dir, "two.txt"), []byte("2"), 0644)

	//listing the contents of a directory
	entries, _ := os.ReadDir(dir)
	fmt.Println("contents of", dir, ":")
	for _, e := range entries {
		fmt.Println(" -", e.Name(), "isDir:", e.IsDir())
	}

	//breaking a path into pieces
	sample := filepath.Join(dir, "reports", "2024", "summary.pdf")
	fmt.Println("Base:", filepath.Base(sample)) // summary.pdf
	fmt.Println("Dir :", filepath.Dir(sample))  // .../reports/2024
	fmt.Println("Ext :", filepath.Ext(sample))  // .pdf

	//turning a relative path into an absolute one
	abs, _ := filepath.Abs("data.txt")
	fmt.Println("absolute path:", abs)

	//walking a whole directory tree - visits every file and folder inside
	fmt.Println("walking", dir, ":")
	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(" ", path)
		return nil
	})

	//removing a directory and everything inside it
	os.RemoveAll(subDir)
	_, err = os.Stat(subDir)
	fmt.Println("sub removed, exists now?", !os.IsNotExist(err))

}
