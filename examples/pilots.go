package main

import (
	"fmt"
	"os"

	"github.com/jameycribbs/hare"
)

type Pilot struct {
	// Required field
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Experience string `json:"experience"`
	Country    string `json:"country"`
}

// Required method
func (pilot *Pilot) SetID(id int) {
	pilot.ID = id
}

// Required method
func (pilot *Pilot) GetID() int {
	return pilot.ID
}

func main() {
	// Open the database
	db, err := hare.OpenDB("data")
	if err != nil {
		fmt.Println("Failed to open database:", err)
		os.Exit(1)
	}

	defer db.Close()

	// Create the pilots table
	pilotsTbl, err := db.CreateTable("pilots")
	if err != nil {
		panic(err)
	}

	// Add a pilot
	pilotRecID, err := pilotsTbl.Create(&Pilot{Name: "Douglas Baeder", Experience: "veteran", Country: "Great Britain"})

	var pilot Pilot

	// Find a pilot
	err = pilotsTbl.Find(pilotRecID, &pilot)
	if err != nil {
		panic(err)
	}

	fmt.Println("ID:", pilot.ID)
	fmt.Println("Name:", pilot.Name)
	fmt.Println("Experience:", pilot.Experience)
	fmt.Println("Country:", pilot.Country)

	// Drop the pilots table
	if err = db.DropTable("pilots"); err != nil {
		panic(err)
	}

}
