package models

import (
	"github.com/jameycribbs/hare"
)

type Host struct {
	// Required field!!!
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (h *Host) GetID() int {
	return h.ID
}

func (h *Host) SetID(id int) {
	h.ID = id
}

func (h *Host) AfterFind(db *hare.Database) error {
	// IMPORTANT!!!  This line of code is necessary in your AfterFind
	//               in order for the Find method to work correctly!
	*h = Host(*h)

	return nil
}

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
