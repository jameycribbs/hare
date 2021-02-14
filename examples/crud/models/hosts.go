package models

import (
	"github.com/jameycribbs/hare"
)

// Host is a record for a MST3K episode comment.
type Host struct {
	// Required field!!!
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetID returns the record id.
// This method is used internally by Hare.
// You need to add this method to each one of
// your models.
func (h *Host) GetID() int {
	return h.ID
}

// SetID takes an id. This method is used
// internally by Hare.
// You need to add this method to each one of
// your models.
func (h *Host) SetID(id int) {
	h.ID = id
}

// AfterFind is a callback that is run by Hare after
// a record is found.
// You need to add this method to each one of your
// models.
func (h *Host) AfterFind(db *hare.Database) error {
	// IMPORTANT!!!  These two lines of code are necessary in your AfterFind
	//               in order for the Find method to work correctly!
	*h = Host(*h)

	return nil
}

// QueryHosts takes a Hare db handle and a query function, and returns
// an array of comments.  If you add this boilerplate method to your model
// you can then write queries using a closure as the query language.
func QueryHosts(db *hare.Database, queryFn func(h Host) bool, limit int) ([]Host, error) {
	var results []Host
	var err error

	ids, err := db.IDs("hosts")
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		h := Host{}

		if err = db.Find("hosts", id, &h); err != nil {
			return nil, err
		}

		if queryFn(h) {
			results = append(results, h)
		}

		if limit != 0 && limit == len(results) {
			break
		}
	}

	return results, err
}
