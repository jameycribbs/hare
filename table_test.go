package hare

import (
	"os"
	"os/exec"
	"reflect"
	"sort"
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
			if filmShouldBe != rec.Film {
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
			if filmShouldBe != rec.Film {
				t.Errorf("want %v; got %v", filmShouldBe, rec.Film)
			}
		},
		func(t *testing.T) {
			//Update...

			db := openTestDB()
			defer db.Close()

			tbl, err := db.GetTable("test")
			if err != nil {
				t.Fatal(err)
			}

			rec := record{}
			err = tbl.Find(3, &rec)
			if err != nil {
				t.Fatal(err)
			}

			rec.Film = rec.Film + " - The Final Cut"
			if err = tbl.Update(&rec); err != nil {
				t.Fatal(err)
			}

			rec = record{}
			err = tbl.Find(3, &rec)
			if err != nil {
				t.Fatal(err)
			}

			filmShouldBe := "The Skydivers - The Final Cut"
			if filmShouldBe != rec.Film {
				t.Errorf("want %v; got %v", filmShouldBe, rec.Film)
			}
		},
		func(t *testing.T) {
			//Destroy...

			db := openTestDB()
			defer db.Close()

			tbl, err := db.GetTable("test")
			if err != nil {
				t.Fatal(err)
			}

			err = tbl.Destroy(3)
			if err != nil {
				t.Fatal(err)
			}

			if err = tbl.Find(3, &record{}); err == nil {
				t.Errorf("want %v; got %v", "Find Error: Record with ID of 3 does not exist!", err)
			}
		},
		func(t *testing.T) {
			//IDs...

			db := openTestDB()
			defer db.Close()

			tbl, err := db.GetTable("test")
			if err != nil {
				t.Fatal(err)
			}

			ids := tbl.IDs()
			sort.Ints(ids)

			idsShouldBe := []int{1, 2, 3, 4}

			if !reflect.DeepEqual(idsShouldBe, ids) {
				t.Errorf("want %v; got %v", idsShouldBe, ids)
			}
		},
		func(t *testing.T) {
			//incrementLastID...

			db := openTestDB()
			defer db.Close()

			tbl, err := db.GetTable("test")
			if err != nil {
				t.Fatal(err)
			}

			id := tbl.incrementLastID()

			idShouldBe := 5

			if idShouldBe != id {
				t.Errorf("want %v; got %v", idShouldBe, id)
			}
		},
		func(t *testing.T) {
			//initIndex...

			db := openTestDB()
			defer db.Close()

			tbl, err := db.GetTable("test")
			if err != nil {
				t.Fatal(err)
			}

			tbl.initIndex()

			indexShouldBe := make(map[int]int64)
			indexShouldBe[1] = 0
			indexShouldBe[2] = 222
			indexShouldBe[3] = 398
			indexShouldBe[4] = 578

			if !reflect.DeepEqual(indexShouldBe, tbl.index) {
				t.Errorf("want %v; got %v", indexShouldBe, tbl.index)
			}
		},
		func(t *testing.T) {
			//initLastID...

			db := openTestDB()
			defer db.Close()

			tbl, err := db.GetTable("test")
			if err != nil {
				t.Fatal(err)
			}

			lastIDShouldBe := 4

			if lastIDShouldBe != tbl.lastID {
				t.Errorf("want %v; got %v", lastIDShouldBe, tbl.lastID)
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
