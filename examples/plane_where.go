package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jameycribbs/hare"
)

type Plane struct {
	ID         int
	Name       string
	Speed      int
	Range      int
	EngineType string
	Country    string
	PlaneType  string
}

func newPlane(rawRec map[string]interface{}) Plane {
	return Plane{
		ID:         int(rawRec["id"].(float64)),
		Name:       rawRec["name"].(string),
		Speed:      int(rawRec["speed"].(float64)),
		Range:      int(rawRec["range"].(float64)),
		EngineType: rawRec["enginetype"].(string),
		Country:    rawRec["country"].(string),
		PlaneType:  rawRec["planetype"].(string),
	}
}

func (plane *Plane) print() {
	fmt.Println("ID:", plane.ID,
		"Name:", plane.Name,
		" | Speed:", plane.Speed,
		" | Range:", plane.Range,
		" | EngineType:", plane.EngineType,
		" | Country:", plane.Country,
		" | PlaneType:", plane.PlaneType,
	)
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

	fmt.Println("Querying for USA...")

	results, err := planesTbl.Where("country == 'USA'")
	if err != nil {
		panic(err)
	}

	for _, result := range results {
		plane := newPlane(result)

		plane.print()
	}

	fmt.Println("Querying for speed less than 375...")

	results, err = planesTbl.Where("speed < 375")
	if err != nil {
		panic(err)
	}

	for _, result := range results {
		plane := newPlane(result)

		plane.print()
	}
}
