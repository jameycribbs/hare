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

	var plane Plane

	fmt.Println("Finding ME-109...")

	if err = planesTbl.Find(2, &plane); err != nil {
		panic(err)
	}

	fmt.Println("ID:", plane.ID)
	fmt.Println("Name:", plane.Name)
	fmt.Println("Speed:", plane.Speed)
	fmt.Println("Range:", plane.Range)
	fmt.Println("EngineType:", plane.EngineType)
	fmt.Println("Country:", plane.Country)
	fmt.Println("PlaneType:", plane.PlaneType)
}
