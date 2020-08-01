package models

import (
	"time"

	"github.com/jameycribbs/hare"
)

type Episodes struct {
	*hare.Table
}

func (episodes *Episodes) Query(queryFn func(episode Episode) bool, limit int) ([]Episode, error) {
	var results []Episode
	var err error

	for _, id := range episodes.Table.IDs() {
		episode := Episode{}

		if err = episodes.Table.Find(id, &episode); err != nil {
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

type Episode struct {
	// Required field!!!
	ID               int       `json:"id"`
	Season           int       `json:"season"`
	Episode          int       `json:"episode"`
	Film             string    `json:"film"`
	Shorts           []string  `json:"shorts"`
	YearFilmReleased int       `json:"year_film_released"`
	DateEpisodeAired time.Time `json:"date_episode_aired"`
	Host             string    `json:"host"`
}

func (episode *Episode) GetID() int {
	return episode.ID
}

func (episode *Episode) SetID(id int) {
	episode.ID = id
}

func (episode *Episode) AfterFind() {
	*episode = Episode(*episode)
}
