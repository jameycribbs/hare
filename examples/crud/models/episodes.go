package models

import (
	"time"

	"github.com/jameycribbs/hare"
)

type Episode struct {
	// Required field!!!
	ID               int       `json:"id"`
	Season           int       `json:"season"`
	Episode          int       `json:"episode"`
	Film             string    `json:"film"`
	Shorts           []string  `json:"shorts"`
	YearFilmReleased int       `json:"year_film_released"`
	DateEpisodeAired time.Time `json:"date_episode_aired"`
	HostID           int       `json:"host_id"`
	Host
	Comments []Comment
}

func (e *Episode) GetID() int {
	return e.ID
}

func (e *Episode) SetID(id int) {
	e.ID = id
}

func (e *Episode) AfterFind(db *hare.Database) error {
	// IMPORTANT!!!  This line of code is necessary in your AfterFind
	//               in order for the Find method to work correctly!
	*e = Episode(*e)

	// This is an example of how you can do Rails-like associations.
	// When an episode is found, this code will run and lookup the
	// associated host record then populate the embedded Host
	// struct.
	h := Host{}
	err := db.Find("hosts", e.HostID, &h)
	if err != nil {
		return err
	} else {
		e.Host = h
	}

	// This is an example of how you can do a Rails-like "has_many"
	// association.  This will run a query on the comments table and
	// populate the episode's Comments embedded struct with child
	// comment records.
	e.Comments, err = QueryComments(db, func(c Comment) bool {
		return c.EpisodeID == e.ID
	}, 0)
	if err != nil {
		return err
	}

	return nil
}

func QueryEpisodes(db *hare.Database, queryFn func(e Episode) bool, limit int) ([]Episode, error) {
	var results []Episode
	var err error

	ids, err := db.IDs("episodes")
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		e := Episode{}

		if err = db.Find("episodes", id, &e); err != nil {
			return nil, err
		}

		if queryFn(e) {
			results = append(results, e)
		}

		if limit != 0 && limit == len(results) {
			break
		}
	}

	return results, err
}
