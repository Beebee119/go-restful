package db

import (
	"database/sql"
	"log"
)

func Initialize(dbDriver *sql.DB) {
	// Create tble train
	statement, err := dbDriver.Prepare(train)
	if err != nil {
		log.Println(err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Println("Error in creating table train ", err)
	} else {
		log.Println("Successfully created table train")
	}

	// Create tble station
	statement, err = dbDriver.Prepare(station)
	if err != nil {
		log.Println(err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Println("Error in creating table station ", err)
	} else {
		log.Println("Successfully created table station")
	}

	// Create tble schedule
	statement, err = dbDriver.Prepare(schedule)
	if err != nil {
		log.Println(err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Println("Error in creating table schedule ", err)
	} else {
		log.Println("Successfully created table schedule")
	}
}
