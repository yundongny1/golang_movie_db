package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

// func main() {
// 	database, _ := sql.Open("sqlite", "./movie_db.db")
// 	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS imdb_movies (id PRIMARY INTEGER KEY, name TEXT, year INTEGER, rank REAL)")
// 	statement.Exec()

// 	statement, _ = database.Prepare("INSERT INTO imdb_movies (id, name, year, rank) VALUES (?, ?, ?, ?)")
// 	statement.Exec("0", "#28", "2002", "12.5")
// 	rows, _ := database.Query("SELECT * FROM imdb_movies")

// 	var id int
// 	var name string
// 	var year int
// 	var rank float32
// 	for rows.Next() {
// 		rows.Scan(&id, &name, &year, &rank)
// 		fmt.Println(id)
// 		rows.Close()
// 	}
// }

type Movie struct {
	ID   string
	Name string
	Year string
	Rank string
}

// func NewActivities() (*Activities, error) {
// 	db, err := sql.Open("sqlite", file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if _, err := db.Exec(create); err != nil {
// 		return nil, err
// 	}
// 	return &Activities{
// 		db: db,
// 	}, nil
// }

func main() {
	var i int

	fmt.Print("Type 1 to truncate and load the IMDB-movies.csv in to the database. Type 2 to select a sample row from the database.")
	fmt.Scan(&i)

	//if user inputs 1, we truncate and load the imdb_movies
	if i == 1 {
		//open the db
		start := time.Now()
		db, err := sql.Open("sqlite", "./movie_db.db")
		if err != nil {
			log.Fatalf("failed to open db: %v", err)
		}
		defer db.Close()

		//open the csv
		file, err := os.Open("IMDB-movies.csv")

		if err != nil {
			log.Fatalf("failed to open file: %v", err)
		}
		defer file.Close()

		// Create a new CSV reader
		reader := csv.NewReader(file)

		// Skip the header row
		if _, err := reader.Read(); err != nil {
			log.Fatalf("failed to read header row: %v", err)
		}

		//truncate and load
		deleteTableSQL := `DROP TABLE IF EXISTS imdb_movies;`
		_, err = db.Exec(deleteTableSQL)
		if err != nil {
			log.Fatalf("Failed to delete table: %v", err)
		}

		createTableSQL := `CREATE TABLE IF NOT EXISTS imdb_movies (
			"id" INTEGER NOT NULL PRIMARY KEY,
			"name" TEXT,
			"year" INTEGER,
			"rank" REAL
		);`

		_, err = db.Exec(createTableSQL)
		if err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}

		fmt.Println("\nTable created successfully! \nInserting rows now.")
		// Prepare a slice to hold the records
		var movie []Movie
		var error_count int

		// Process each line of the CSV and insert
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				//Log the error, skip the problematic line, and increment error count
				//log.Printf("Error reading record: %v. Skipping this line.", err)
				//error_records = append(err)
				error_count++
				continue
			}

			// Convert the fields to the appropriate types
			id := record[0]
			// if err != nil {
			// 	log.Printf("Error converting ID: %v. Skipping this line.", err)
			// 	continue
			// }

			name := record[1]
			// if err != nil {
			// 	log.Printf("Error converting Age: %v. Skipping this line.", err)
			// 	continue
			// }

			year := record[2]
			// if err != nil {
			// 	log.Printf("Error converting Age: %v. Skipping this line.", err)
			// 	continue
			// }

			ranks := record[3]
			// if err != nil {
			// 	log.Printf("Error converting Age: %v. Skipping this line.", err)
			// 	continue
			// }

			// Create a new Person struct with the data
			movies := Movie{
				ID:   id,
				Name: name,
				Year: year,
				Rank: ranks,
			}

			// Append the struct to the slice

			movie = append(movie, movies)
			stmt, err := db.Prepare("INSERT INTO imdb_movies (id, name, year, rank) values(?,?,?,?)")
			if err != nil {
				log.Fatalf("failed to insert into imdb_movies: %v", err)
			}

			_, err = stmt.Exec(movies.ID, movies.Name, movies.Year, movies.Rank)
			if err != nil {
				log.Fatalf("failed to insert values: %v", err)
			}
		}
		fmt.Printf("Loaded %v, skipped %v records", len(movie), error_count)
		elapsed := time.Since(start)
		fmt.Printf("Insert execution time: %s\n", elapsed)

	} else if i == 2 {
		var j int
		fmt.Println("Enter the movie ID you want to retrieve.")
		fmt.Scan(&j)

		sqlStatement := `SELECT * FROM imdb_movies WHERE id=$1;`
		var id string
		var title string
		var year string
		var rank string

		db, err := sql.Open("sqlite", "./movie_db.db")
		if err != nil {
			log.Fatalf("failed to open db: %v", err)
		}
		defer db.Close()

		row := db.QueryRow(sqlStatement, j)
		switch err := row.Scan(&id, &title, &year, &rank); err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
		case nil:
			fmt.Println(id, title, year, rank)
		default:
			panic(err)
		}
	}
}
