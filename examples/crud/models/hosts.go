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

func (h *Host) AfterFind(db *hare.Database) {
	*h = Host(*h)
}

func QueryHosts(db *hare.Database, queryFn func(host Host) bool, limit int) ([]Host, error) {
	var results []Host
	var err error

	ids, err := db.IDs("mst3k_hosts")
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		host := Host{}

		if err = db.Find("mst3k_hosts", id, &host); err != nil {
			return nil, err
		}

		if queryFn(host) {
			results = append(results, host)
		}

		if limit != 0 && limit == len(results) {
			break
		}
	}

	return results, err
}
