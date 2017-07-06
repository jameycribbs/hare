package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jameycribbs/hare"
)

type Plane struct {
	// Required field
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Speed      int    `json:"speed"`
	Range      int    `json:"range"`
	EngineType string `json:"enginetype"`
	Country    string `json:"country"`
	PlaneType  string `json:"planetype"`
}

// Required method
func (plane *Plane) SetID(id int) {
	plane.ID = id
}

// Required method
func (plane *Plane) GetID() int {
	return plane.ID
}

func main() {
	cmd := exec.Command("cp", "data/planes_default.json", "data/planes.json")

	if err := cmd.Run(); err != nil {
		panic(err)
	}

	// Open the database
	db, err := hare.OpenDB("data")
	if err != nil {
		fmt.Println("Failed to open database:", err)
		os.Exit(1)
	}

	defer db.Close()

	// Grab a reference to the planes table
	planesTbl, err := db.GetTable("planes")
	if err != nil {
		panic(err)
	}

	// Find
	err = planesTbl.ForEachID(func(recID int) error {
		var plane Plane

		if err = planesTbl.Find(recID, &plane); err != nil {
			panic(err)
		}

		if plane.Speed < 360 {
			fmt.Println("Plane with speed less than 360:", plane)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	var foundSpitfire Plane

	// Find Spitfire...
	err = planesTbl.ForEachID(func(recID int) error {
		var plane Plane

		if err = planesTbl.Find(recID, &plane); err != nil {
			panic(err)
		}

		if plane.Name == "Spitfire I" {
			foundSpitfire = plane

			// If you want to exit the ForEach early, say, for example, because you
			// found the record you were looking for, you need to do this inside
			// your closure:
			return hare.ForEachIDBreak{}
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	if foundSpitfire.ID == 0 {
		fmt.Println("Spitfire not found!")
	} else {
		fmt.Println("Found Spitfire:", foundSpitfire)
	}
}
