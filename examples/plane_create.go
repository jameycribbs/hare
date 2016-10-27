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

func (plane *Plane) Populate(recMap map[string]interface{}) {
	plane.ID = int(recMap["id"].(float64))
	plane.Name = recMap["name"].(string)
	plane.Speed = int(recMap["speed"].(float64))
	plane.Range = int(recMap["range"].(float64))
	plane.EngineType = recMap["enginetype"].(string)
	plane.Country = recMap["country"].(string)
	plane.PlaneType = recMap["planetype"].(string)
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

	id, err := planesTbl.Create(&Plane{Name: "FW-190", Speed: 389, Range: 555, EngineType: "radial", Country: "Germany", PlaneType: "fighter"})

	if err != nil {
		panic(err)
	}

	var plane Plane

	if err = planesTbl.Find(id, &plane); err != nil {
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
