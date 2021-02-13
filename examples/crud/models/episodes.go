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
}

func (episode *Episode) GetID() int {
	return episode.ID
}

func (episode *Episode) SetID(id int) {
	episode.ID = id
}

func (episode *Episode) AfterFind(db *hare.Database) {
	*episode = Episode(*episode)

	// This is an example of how you can do Rails-like associations.
	// When an episode is found, this code will run and lookup the
	// associated host record then populate the embedded Host
	// struct.
	host := Host{}
	err := db.Find("mst3k_hosts", episode.HostID, &host)

	if err == nil {
		episode.Host = host
	}
}

func QueryEpisodes(db *hare.Database, queryFn func(episode Episode) bool, limit int) ([]Episode, error) {
	var results []Episode
	var err error

	ids, err := db.IDs("mst3k_episodes")
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		episode := Episode{}

		if err = db.Find("mst3k_episodes", id, &episode); err != nil {
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
