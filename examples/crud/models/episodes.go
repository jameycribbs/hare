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

func (episode *Episode) GetID() int {
	return episode.ID
}

func (episode *Episode) SetID(id int) {
	episode.ID = id
}

func (episode *Episode) AfterFind(db *hare.Database) error {
	// IMPORTANT!!!  This line of code is necessary in your AfterFind
	//               in order for the Find method to work correctly!
	*episode = Episode(*episode)

	// This is an example of how you can do Rails-like associations.
	// When an episode is found, this code will run and lookup the
	// associated host record then populate the embedded Host
	// struct.
	host := Host{}
	err := db.Find("hosts", episode.HostID, &host)
	if err != nil {
		return err
	} else {
		episode.Host = host
	}

	// This is an example of how you can do a Rails-like "has_many"
	// association.  This will run a query on the comments table and
	// populate the episode's Comments embedded struct with child
	// comment records.
	episode.Comments, err = QueryComments(db, func(c Comment) bool {
		return c.EpisodeID == episode.ID
	}, 0)
	if err != nil {
		return err
	}

	return nil
}

func QueryEpisodes(db *hare.Database, queryFn func(episode Episode) bool, limit int) ([]Episode, error) {
	var results []Episode
	var err error

	ids, err := db.IDs("episodes")
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		episode := Episode{}

		if err = db.Find("episodes", id, &episode); err != nil {
			return nil, err
		}

		if queryFn(episode) {
			results = append(results, episode)
		}

		if limit != 0 && limit == len(results) {
			break
		}
	}

	return results, err
}
