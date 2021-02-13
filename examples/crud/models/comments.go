package models

import (
	"github.com/jameycribbs/hare"
)

type Comment struct {
	// Required field!!!
	ID        int    `json:"id"`
	EpisodeID int    `json:"episode_id"`
	Text      string `json:"text"`
}

func (c *Comment) GetID() int {
	return c.ID
}

func (c *Comment) SetID(id int) {
	c.ID = id
}

func (c *Comment) AfterFind(db *hare.Database) error {
	// IMPORTANT!!!  This line of code is necessary in your AfterFind
	//               in order for the Find method to work correctly!
	*c = Comment(*c)

	return nil
}

func QueryComments(db *hare.Database, queryFn func(c Comment) bool, limit int) ([]Comment, error) {
	var results []Comment
	var err error

	ids, err := db.IDs("comments")
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		c := Comment{}

		if err = db.Find("comments", id, &c); err != nil {
			return nil, err
		}

		if queryFn(c) {
			results = append(results, c)
		}

		if limit != 0 && limit == len(results) {
			break
		}
	}

	return results, err
}
