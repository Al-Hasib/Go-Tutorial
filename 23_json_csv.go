package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// struct tags after the type tell encoding/json how to name the field
// in JSON, and whether to skip it when it's empty.
type Person struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Email   string   `json:"email,omitempty"` // omitted from JSON if ""
	private string   // lowercase = unexported, json package can't see it at all
	Hobbies []string `json:"hobbies"`
}

func main() {

	// everything below happens inside a temp folder so this lesson
	// doesn't leave files behind in your project folder
	dir, err := os.MkdirTemp("", "go-lesson-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	// ---------- JSON: MARSHAL (Go value -> JSON bytes) ----------

	p := Person{Name: "Hasib", Age: 24, Hobbies: []string{"coding", "chess"}}

	//json.Marshal turns a Go value into compact JSON bytes.
	//note "email" is missing (omitempty) and "private" never appears at all.
	data, err := json.Marshal(p)
	if err != nil {
		fmt.Println("marshal error:", err)
	}
	fmt.Println("compact JSON:", string(data))

	//json.MarshalIndent adds indentation - nicer for humans/logs/files
	pretty, _ := json.MarshalIndent(p, "", "  ")
	fmt.Println("pretty JSON:\n" + string(pretty))

	//marshaling a slice of structs -> a JSON array
	people := []Person{
		{Name: "Alice", Age: 30, Email: "alice@example.com"},
		{Name: "Bob", Age: 25},
	}
	listJSON, _ := json.MarshalIndent(people, "", "  ")
	fmt.Println("slice of structs JSON:\n" + string(listJSON))

	// ---------- JSON: UNMARSHAL (JSON bytes -> Go value) ----------

	jsonInput := []byte(`{"name":"Rana","age":28,"email":"rana@example.com","hobbies":["reading"]}`)

	//unmarshal into a struct we already know the shape of - the common case
	var decoded Person
	err = json.Unmarshal(jsonInput, &decoded) // needs a pointer, it fills the struct in place
	if err != nil {
		fmt.Println("unmarshal error:", err)
	}
	fmt.Printf("decoded struct: %+v\n", decoded)

	//unmarshal into map[string]interface{} when you don't know (or don't
	//care about) the exact shape - every value comes back as `any`
	var asMap map[string]interface{}
	json.Unmarshal(jsonInput, &asMap)
	fmt.Println("decoded as map:", asMap)
	fmt.Println("  age from map (float64 by default):", asMap["age"])

	// ---------- JSON: READING/WRITING FILES ----------

	jsonPath := filepath.Join(dir, "people.json")

	//write JSON straight to a file: marshal, then os.WriteFile
	if b, err := json.MarshalIndent(people, "", "  "); err == nil {
		os.WriteFile(jsonPath, b, 0644)
	}
	fmt.Println("wrote JSON file:", jsonPath)

	//read it back: os.ReadFile, then unmarshal into a matching slice
	raw, _ := os.ReadFile(jsonPath)
	var loaded []Person
	json.Unmarshal(raw, &loaded)
	fmt.Println("loaded back from file:", loaded)

	//json.NewEncoder/NewDecoder work directly on an *os.File - handy when
	//the data is big and you don't want the whole thing in memory at once
	if f, err := os.Create(jsonPath); err == nil {
		json.NewEncoder(f).Encode(people) // encoder writes straight to the file
		f.Close()
	}
	if f, err := os.Open(jsonPath); err == nil {
		var streamed []Person
		json.NewDecoder(f).Decode(&streamed) // decoder reads straight from the file
		f.Close()
		fmt.Println("streamed via encoder/decoder:", streamed)
	}

	// ---------- CSV: WRITING ----------

	csvPath := filepath.Join(dir, "people.csv")

	csvFile, err := os.Create(csvPath)
	if err != nil {
		fmt.Println("create csv error:", err)
	}
	writer := csv.NewWriter(csvFile)

	//header row first, then one []string per record
	writer.Write([]string{"name", "age", "email"})
	for _, per := range people {
		writer.Write([]string{per.Name, fmt.Sprint(per.Age), per.Email})
	}
	writer.Flush() // buffered - nothing hits the file until Flush (or Close)
	csvFile.Close()
	fmt.Println("wrote CSV file:", csvPath)

	// ---------- CSV: READING ----------

	readFile, err := os.Open(csvPath)
	if err != nil {
		fmt.Println("open csv error:", err)
	}
	defer readFile.Close()

	reader := csv.NewReader(readFile)

	//ReadAll loads every record into memory as [][]string - simplest option
	//for small/medium files
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("read csv error:", err)
	}
	fmt.Println("all CSV records:", records)

	header, rows := records[0], records[1:]
	fmt.Println("header:", header)
	for _, row := range rows {
		fmt.Printf("  name=%s age=%s email=%s\n", row[0], row[1], row[2])
	}

	//reading record-by-record with Read() - better for big files, since it
	//doesn't load the whole thing into memory at once
	readFile2, _ := os.Open(csvPath)
	defer readFile2.Close()
	reader2 := csv.NewReader(readFile2)
	fmt.Println("reading record by record:")
	for {
		record, err := reader2.Read()
		if err != nil {
			break // io.EOF once every record has been read
		}
		fmt.Println(" ", record)
	}
}
