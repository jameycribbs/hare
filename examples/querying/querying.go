package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/jameycribbs/hare"
)

// IMPORTANT!
// Your record's struct MUST have an "ID" field and it must
// also implement the 3 methods below in order for it to
// satisfy the hare.Record interface.

type episodeRec struct {
	ID               int       `json:"id"`
	Season           int       `json:"season"`
	Episode          int       `json:"episode"`
	Film             string    `json:"film"`
	Shorts           []string  `json:"shorts"`
	YearFilmReleased int       `json:"year_film_released"`
	DateEpisodeAired time.Time `json:"date_episode_aired"`
	Host             string    `json:"host"`
}

func (e *episodeRec) GetID() int {
	return e.ID
}

func (e *episodeRec) SetID(id int) {
	e.ID = id
}

func (e *episodeRec) AfterFind() {
	*e = episodeRec(*e)
}

// To allow easy querying of a Hare table, you need to create a new
// struct type with an embedded hare.Table and implement a query
// method for it.

type episodesMdl struct {
	*hare.Table
}

func (eps *episodesMdl) query(queryFn func(rec episodeRec) bool, limit int) ([]episodeRec, error) {
	var results []episodeRec
	var err error

	for _, id := range eps.Table.IDs() {
		rec := episodeRec{}

		if err = eps.Table.Find(id, &rec); err != nil {
			return nil, err
		}

		if queryFn(rec) {
			results = append(results, rec)
		}

		if limit != 0 && limit == len(results) {
			break
		}
	}

	return results, err
}

func main() {
	setupExampleDB()

	var episodes episodesMdl

	// Open the database and return a handle to it.
	db, err := hare.OpenDB("../example_data")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Let's grab a handle to the MST3K Episodes table
	// so we can play with it.
	episodes.Table, err = db.GetTable("mst3k_episodes")
	if err != nil {
		panic(err)
	}

	// Now we will run a query for episodes that Joel hosted.
	// Notice that we are actually passing the query expression as an
	// anonymous function.  This allows us to use power of Go itself to
	// query the database.  There is no need to learn and use a DSL.
	results, err := episodes.query(func(r episodeRec) bool {
		return r.Host == "Joel"
	}, 0)
	if err != nil {
		panic(err)
	}

	for _, r := range results {
		fmt.Printf("Joel hosted the season %v episode %v film, '%v'\n", r.Season, r.Episode, r.Film)
	}
}

func setupExampleDB() {
	cmd := exec.Command("cp", "../example_data/mst3k_episodes_default.json", "../example_data/mst3k_episodes.json")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
