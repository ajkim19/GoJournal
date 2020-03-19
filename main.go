package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB
var temps *template.Template
var err error
var dbid int
var dbdate string
var dbentry string
var rows *sql.Rows
var username string = "journal"
var JEntries = []JEntry{}

type JEntry struct {
	Date  string `json:"date"`
	Entry string `json:"entry"`
}

func init() {
	dataSourceName := fmt.Sprintf("%s.db", username)

	// Makes a handle for the database journal
	database, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	// Creates the table journal_entries if it has been dropped
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS journal_entries (id INTEGER PRIMARY KEY, date TEXT, entry TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()

	rows, err = database.Query("SELECT * FROM journal_entries")
	if err != nil {
		log.Fatal(err)
	}

	var dateExists bool

	for rows.Next() {
		rows.Scan(&dbid, &dbdate, &dbentry)
		if dbdate == "2020-01-06" {
			dateExists = true
		}
	}

	// Adds an entry to journal_entries if it is empty
	if dateExists == false {
		statement, err := database.Prepare(`
			INSERT INTO journal_entries (date, entry)
			VALUES
				("2020-01-06", "Welcome to Revature! Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."),
				("2020-01-15", "It's Sunday! Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."),
				("2020-01-16", "It's Monday! Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."),
				("2020-01-17", "It's Tuesday! Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."),
				("2020-01-18", "It's Wednesday! Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")`)
		if err != nil {
			log.Fatal(err)
		}
		statement.Exec()
	}
}

func main() {
	temps, err = template.ParseGlob("main.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rows, err = database.Query("SELECT * FROM journal_entries ORDER BY date DESC")
		if err != nil {
			log.Fatal(err)
		}
		for rows.Next() {
			rows.Scan(&dbid, &dbdate, &dbentry)
			JEntries = append(JEntries, JEntry{Date: dbdate, Entry: dbentry})
		}

		err = temps.Execute(w, JEntries)
		if err != nil {
			log.Fatal(err)
		}
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
	// for _, entry := range JEntries {
	// 	fmt.Println(entry.Date, entry.Entry)
	// 	fmt.Println()
	// }
}
