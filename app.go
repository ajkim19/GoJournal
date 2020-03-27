package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB
var tmpl *template.Template
var err error
var dbid int
var dbdate string
var dbentry string
var rows *sql.Rows
var username string = "journal"
var JEntries []JEntry
var journalDate string
var journalEntry string

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
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	tmpl, err = template.ParseFiles("index.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err = r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		JEntries = []JEntry{}

		rows, err = database.Query("SELECT * FROM journal_entries ORDER BY date DESC")
		if err != nil {
			log.Fatal(err)
		}
		for rows.Next() {
			rows.Scan(&dbid, &dbdate, &dbentry)
			JEntries = append(JEntries, JEntry{Date: dbdate, Entry: dbentry})
		}

		err = tmpl.Execute(w, JEntries)
		if err != nil {
			log.Fatal(err)
		}

		/*
			///////////////////////////////////////////////
			Inputting entries
			///////////////////////////////////////////////
		*/
		// Editing an existing entry
		if r.FormValue("edit-entry") != "" {
			journalDate = r.FormValue("edit-entry-date")
			journalEntry = r.FormValue("edit-entry")

			//fmt.Println(journalEntry)

			updateEntry(database, journalDate, journalEntry)

			// Adding an entry
		} else if r.FormValue("entry") != "" {
			journalDate = r.FormValue("date")
			journalEntry = r.FormValue("entry")

			rows, err = database.Query(`SELECT * FROM journal_entries WHERE date = ?`, journalDate)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()

			dateExists := false

			for rows.Next() {
				err := rows.Scan(&dbid, &dbdate, &dbentry)
				if err != nil {
					log.Fatal(err)
				}
				if journalDate == dbdate {
					dateExists = true
				}
			}

			// If the date of the entry already exists, the entry will be added	to
			// the preexisting entry after a new line.
			if dateExists {
				rows, err = database.Query("SELECT * FROM journal_entries WHERE date = ?", journalDate)
				if err != nil {
					log.Fatal(err)
				}
				for rows.Next() {
					err := rows.Scan(&dbid, &dbdate, &dbentry)
					if err != nil {
						log.Fatal(err)
					}
					journalEntry = fmt.Sprint(dbentry + "\n\n" + journalEntry)
				}

				updateEntry(database, journalDate, journalEntry)

			} else {
				statement, err := database.Prepare("INSERT INTO journal_entries (date, entry) VALUES (?, ?)")
				if err != nil {
					log.Fatal(err)
				}
				defer statement.Close()
				statement.Exec(journalDate, journalEntry)
			}
		}
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// checkDateFormat checks to see if the inputted date is in the correct format
func checkDateFormat(date string) string {
	matched, err := regexp.MatchString(`((19|20)[0-9][0-9])[- /.](0[1-9]|1[012])[- /.]([012][0-9]|3[01])`, journalDate)
	if err != nil {
		log.Fatal(err)
	}
	if matched == false {
		journalDate = string(time.Now().Format("2006-01-02"))
	}
	return journalDate
}

// updateEntry edits an entry in the journal database
func updateEntry(db *sql.DB, jd string, je string) {
	statement, err := db.Prepare("UPDATE journal_entries SET entry = ? WHERE date = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()
	statement.Exec(je, jd)
}
