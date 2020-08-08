package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/jameycribbs/hare"
	"github.com/jameycribbs/hare/datastores/disk"
	"github.com/jameycribbs/hare/examples/crud/models"
)

func main() {
	ds, err := disk.New("./data", ".json")
	if err != nil {
		panic(err)
	}

	db, err := hare.New(ds)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//----- CREATE -----

	recID, err := db.Insert("mst3k_episodes", &models.Episode{
		Season:           6,
		Episode:          19,
		Film:             "Red Zone Cuba",
		Shorts:           []string{"Speech:  Platform, Posture, and Appearance"},
		YearFilmReleased: 1966,
		DateEpisodeAired: time.Date(1994, 12, 17, 0, 0, 0, 0, time.UTC),
		Host:             "Mike",
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("New record id is:", recID)

	//----- READ -----

	rec := models.Episode{}

	err = db.Find("mst3k_episodes", 4, &rec)
	if err != nil {
		panic(err)
	}

	fmt.Println("Found record is:", rec.Film)

	//----- UPDATE -----

	rec.Film = "The Skydivers - The Final Cut"
	if err = db.Update("mst3k_episodes", &rec); err != nil {
		panic(err)
	}

	//----- DELETE -----

	err = db.Delete("mst3k_episodes", 2)
	if err != nil {
		panic(err)
	}

	//----- QUERYING -----

	results, err := models.QueryEpisodes(db, func(r models.Episode) bool {
		return r.Host == "Joel"
	}, 0)
	if err != nil {
		panic(err)
	}

	for _, r := range results {
		fmt.Printf("Joel hosted the season %v episode %v film, '%v'\n", r.Season, r.Episode, r.Film)
	}
}

func init() {
	cmd := exec.Command("cp", "./data/mst3k_episodes_default.txt", "./data/mst3k_episodes.json")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
