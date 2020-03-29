package journal

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"time"
)

var dbid int
var dbdate string
var dbentry string

// CheckDateFormat checks to see if the inputted date is in the correct format
func CheckDateFormat(jd string) string {
	matched, err := regexp.MatchString(`((19|20)[0-9][0-9])[- /.](0[1-9]|1[012])[- /.]([012][0-9]|3[01])`, jd)
	if err != nil {
		log.Fatal(err)
	}
	if matched == false {
		jd = string(time.Now().Format("2006-01-02"))
	}
	return jd
}

// AddEntry adds an entry to the journal database
func AddEntry(db *sql.DB, jd string, je string) {
	rows, err := db.Query(`SELECT * FROM journal_entries WHERE date = ?`, jd)
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
		if jd == dbdate {
			dateExists = true
		}
	}

	// If the date of the entry already exists, the entry will be added	to
	// the preexisting entry after a new line.
	if dateExists {
		rows, err = db.Query("SELECT * FROM journal_entries WHERE date = ?", jd)
		if err != nil {
			log.Fatal(err)
		}
		for rows.Next() {
			err := rows.Scan(&dbid, &dbdate, &dbentry)
			if err != nil {
				log.Fatal(err)
			}
			je = fmt.Sprint(dbentry + "\n\n" + je)
		}

		EditEntry(db, jd, je)

	} else {
		statement, err := db.Prepare("INSERT INTO journal_entries (date, entry) VALUES (?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer statement.Close()
		statement.Exec(jd, je)
	}
}

// EditEntry edits an entry in the journal db
func EditEntry(db *sql.DB, jd string, je string) {
	statement, err := db.Prepare("UPDATE journal_entries SET entry = ? WHERE date = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()
	statement.Exec(je, jd)
}

// DeleteEntry deletes an entry in the journal db
func DeleteEntry(db *sql.DB, jd string) {
	statement, err := db.Prepare("DELETE FROM journal_entries WHERE date = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()
	statement.Exec(jd)
}
