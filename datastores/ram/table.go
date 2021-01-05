package ram

type table struct {
	records map[int][]byte
}

func newTable() *table {
	var t table

	t.records = make(map[int][]byte)

	return &t
}

func (t *table) deleteRec(id int) error {
	if !t.recExists(id) {
		return ErrNoRecord
	}

	delete(t.records, id)

	return nil
}

func (t *table) getLastID() int {
	var lastID int

	for id := range t.records {
		if id > lastID {
			lastID = id
		}
	}

	return lastID
}

func (t *table) ids() []int {
	ids := make([]int, len(t.records))

	i := 0
	for id := range t.records {
		ids[i] = id
		i++
	}

	return ids
}

func (t *table) readRec(id int) ([]byte, error) {
	rec, ok := t.records[id]
	if !ok {
		return nil, ErrNoRecord
	}

	return rec, nil
}

func (t *table) recExists(id int) bool {
	_, ok := t.records[id]

	return ok
}

func (t *table) writeRec(id int, rec []byte) {
	t.records[id] = rec
}
