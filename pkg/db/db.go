package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/glebarez/go-sqlite"
	"github.com/google/uuid"
)

var (
	dbPath = os.Getenv("DB_PATH")
)

type Record struct {
	Url                string
	ReplacementSummary string
}

type ChangeRecord struct {
	Id                 string
	Url                string
	ReplacementSummary string
}

func init() {
	if dbPath == "" {
		dbPath = "./data/calendars.db"
	}
}

// TODO: write a class that stores dbPath
func InitDB() error {
	log.Printf("Using db: %s", dbPath)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	stmt := "CREATE TABLE IF NOT EXISTS calendars(id TEXT PRIMARY KEY, url TEXT, replacementSummary TEXT);"
	_, err = db.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func ReadRecord(id string) (Record, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return Record{}, err
	}
	defer db.Close()

	stmt, err := db.Prepare("select url, replacementSummary from calendars where id = ?")
	if err != nil {
		return Record{}, err
	}
	defer stmt.Close()

	var url, replacementSummary string
	err = stmt.QueryRow(id).Scan(&url, &replacementSummary)
	if err == sql.ErrNoRows {
		return Record{}, err
	} else if err != nil {
		return Record{}, err
	}

	return Record{url, replacementSummary}, nil

}

func UpdateRecord(record ChangeRecord) (string, error) {
	fmt.Printf("Updating: %#v\n", record)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Printf("Unable to open sqlite db")
		return "", err
	}
	defer db.Close()

	if record.Id == "" {
		return "", fmt.Errorf("No id passed!")
	} else if record.Url == "" && record.ReplacementSummary == "" {
		return "", fmt.Errorf("Need URL or Summary!")
	}

	// I imagine there's a nicer way to write this :)
	if record.Url != "" && record.ReplacementSummary != "" {
		stmt := "UPDATE calendars set url = ?, replacementSummary = ? where id = ?;"
		_, err = db.Exec(stmt, record.Url, record.ReplacementSummary, record.Id)
		if err != nil {
			log.Printf("Unable to update record in calendars table")
			return "", err
		}
		return "", nil
	} else if record.ReplacementSummary != "" {
		stmt := "UPDATE calendars set replacementSummary = ? where id = ?;"
		_, err = db.Exec(stmt, record.ReplacementSummary, record.Id)
		if err != nil {
			log.Printf("Unable to update summary in calendars table")
			return "", err
		}
		return "", nil
	} else if record.Url != "" {
		stmt := "UPDATE calendars set url = ? where id = ?;"
		_, err = db.Exec(stmt, record.Url, record.Id)
		if err != nil {
			log.Printf("Unable to update url in calendars table")
			return "", err
		}
		return "", nil
	}
	return "", fmt.Errorf("Invalid record! This shouldn't be possible. %#v", record)
}

func WriteRecord(record Record) (string, error) {
	id := uuid.New().String()
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Printf("Unable to open sqlite db")
		return "", err
	}
	defer db.Close()

	stmt := "INSERT INTO calendars(id, url, replacementSummary) VALUES(?, ?, ?);"
	_, err = db.Exec(stmt, id, record.Url, record.ReplacementSummary)
	if err != nil {
		log.Printf("Unable to insert record into calendars table")
		return "", err
	}

	return id, nil

}
