package main

import (
	"fmt"
	"os"

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
	db, err := hare.OpenDB("data")
	if err != nil {
		fmt.Println("Failed to open database:", err)
		os.Exit(1)
	}

	defer db.Close()

	planesTbl, err := db.GetTable("planes")
	if err != nil {
		panic(err)
	}

	err = planesTbl.Update(&Plane{ID: 3, Name: "Spitfire", Speed: 366, Range: 742, EngineType: "inline", Country: "Great Britain", PlaneType: "fighter"})
	err = planesTbl.Update(&Plane{ID: 3, Name: "Spitfire XVI", Speed: 366, Range: 742, EngineType: "inline", Country: "Great Britain", PlaneType: "fighter"})
	err = planesTbl.Update(&Plane{ID: 1, Name: "P-51D", Speed: 403, Range: 1650, EngineType: "inline", Country: "USA", PlaneType: "fighter"})

	hurricaneID, err := planesTbl.Create(&Plane{Name: "Hurricane", Speed: 333, Range: 714, EngineType: "inline", Country: "Great Britain", PlaneType: "fighter"})

	if err != nil {
		panic(err)
	}

	fmt.Println("Hurricane Rec ID:", hurricaneID)

	_, err = planesTbl.Create(&Plane{Name: "FW-190", Speed: 389, Range: 555, EngineType: "radial", Country: "Germany", PlaneType: "fighter"})

	if err != nil {
		panic(err)
	}

	err = planesTbl.Destroy(hurricaneID)
	if err != nil {
		panic(err)
	}

	var plane Plane

	err = planesTbl.Find(4, &plane)
	if err != nil {
		panic(err)
	}

	fmt.Println("Name:", plane.Name)
	fmt.Println("Speed:", plane.Speed)
	fmt.Println("Range:", plane.Range)
	fmt.Println("EngineType:", plane.EngineType)
	fmt.Println("Country:", plane.Country)
	fmt.Println("PlaneType:", plane.PlaneType)

	// All German planes...
	err = planesTbl.ForEachID(func(recID int) error {
		if err = planesTbl.Find(recID, &plane); err != nil {
			return err
		}

		if plane.Country == "Germany" {
			fmt.Println("German Plane:", plane.Name)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Find all German planes failed:", err)
	}

	pilotsTbl, err := db.CreateTable("pilots")
	if err != nil {
		panic(err)
	}

	pilotRecID, err := pilotsTbl.Create(&Pilot{Name: "Douglas Baeder", Experience: "veteran", Country: "Great Britain"})

	var pilot Pilot

	err = pilotsTbl.Find(pilotRecID, &pilot)
	if err != nil {
		panic(err)
	}

	fmt.Println("ID:", pilot.ID)
	fmt.Println("Name:", pilot.Name)
	fmt.Println("Experience:", pilot.Experience)
	fmt.Println("Country:", pilot.Country)

	if err = db.DropTable("pilots"); err != nil {
		panic(err)
	}
}
