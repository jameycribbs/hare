package models

import (
	"github.com/jameycribbs/hare"
)

// Comment is a record for a MST3K episode comment.
type Comment struct {
	// Required field!!!
	ID        int    `json:"id"`
	EpisodeID int    `json:"episode_id"`
	Text      string `json:"text"`
}

// GetID returns the record id.
// This method is used internally by Hare.
// You need to add this method to each one of
// your models.
func (c *Comment) GetID() int {
	return c.ID
}

// SetID takes an id. This method is used
// internally by Hare.
// You need to add this method to each one of
// your models.
func (c *Comment) SetID(id int) {
	c.ID = id
}

// AfterFind is a callback that is run by Hare after
// a record is found.
// You need to add this method to each one of your
// models.
func (c *Comment) AfterFind(db *hare.Database) error {
	// IMPORTANT!!!  These two lines of code are necessary in your AfterFind
	//               in order for the Find method to work correctly!
	*c = Comment(*c)

	return nil
}

// QueryComments takes a Hare db handle and a query function, and returns
// an array of comments.  If you add this boilerplate method to your model
// you can then write queries using a closure as the query language.
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
