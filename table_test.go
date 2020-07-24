package hare

import (
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"
)

type record struct {
	ID               int       `json:"id"`
	Season           int       `json:"season"`
	Episode          int       `json:"episode"`
	Film             string    `json:"film"`
	Shorts           []string  `json:"shorts"`
	YearFilmReleased int       `json:"year_film_released"`
	DateEpisodeAired time.Time `json:"date_episode_aired"`
	Host             string    `json:"host"`
}

func (r *record) GetID() int {
	return r.ID
}

func (r *record) SetID(id int) {
	r.ID = id
}

func (r *record) AfterFind() {
	*r = record(*r)
}

func TestAllTableTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//Find...

			db := openTestDB()
			defer db.Close()

			tbl, err := db.GetTable("test")
			if err != nil {
				t.Fatal(err)
			}

			rec := record{}
			err = tbl.Find(3, &rec)
			if err != nil {
				panic(err)
			}

			filmShouldBe := "The Skydivers"
			if rec.Film != filmShouldBe {
				t.Errorf("want %v; got %v", filmShouldBe, rec.Film)
			}
		},
		func(t *testing.T) {
			//Create...

			db := openTestDB()
			defer db.Close()

			tbl, err := db.GetTable("test")
			if err != nil {
				t.Fatal(err)
			}

			recID, err := tbl.Create(
				&record{
					Season:           6,
					Episode:          19,
					Film:             "Red Zone Cuba",
					Shorts:           []string{"Speech:  Platform, Posture, and Appearance"},
					YearFilmReleased: 1966,
					DateEpisodeAired: time.Date(1994, 12, 17, 0, 0, 0, 0, time.UTC),
					Host:             "Mike",
				})
			if err != nil {
				t.Fatal(err)
			}

			rec := record{}
			err = tbl.Find(recID, &rec)
			if err != nil {
				t.Fatal(err)
			}

			filmShouldBe := "Red Zone Cuba"
			if rec.Film != filmShouldBe {
				t.Errorf("want %v; got %v", filmShouldBe, rec.Film)
			}
		},
	}

	for i, fn := range tests {
		tableTestSetup()
		t.Run(strconv.Itoa(i), fn)
		tableTestTeardown()
	}
}

func tableTestSetup() {
	cmd := exec.Command("cp", "test_data/test_default.json", "test_data/test.json")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func tableTestTeardown() {
	if err := os.Remove("test_data/test.json"); err != nil {
		panic(err)
	}
}
