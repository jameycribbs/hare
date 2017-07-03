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

	var spitfire Plane

	if err = planesTbl.Find(3, &spitfire); err != nil {
		panic(err)
	}

	fmt.Println(spitfire.Name, " starts with speed of:", spitfire.Speed)

	spitfire.Name = "Spitfire XI"
	spitfire.Speed = 366

	fmt.Println("Name will be updated to", spitfire.Name, "and speed will be updated to", spitfire.Speed)

	if err = planesTbl.Update(&spitfire); err != nil {
		panic(err)
	}

	if err = planesTbl.Find(3, &spitfire); err != nil {
		panic(err)
	}

	fmt.Println("Name is successfully updated to", spitfire.Name, "and speed is successfully updated to", spitfire.Speed)
}
