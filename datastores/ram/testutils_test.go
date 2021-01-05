package ram

import (
	"testing"
)

func newTestRam(t *testing.T) *Ram {
	s := make(map[string]map[int]string)
	s["contacts"] = seedData()

	ram, err := NewRam(s)
	if err != nil {
		t.Fatal(err)
	}

	return ram
}

func newTestTable(t *testing.T) *table {
	tbl := newTable()

	for k, v := range seedData() {
		tbl.records[k] = []byte(v)
	}

	return tbl
}

func seedData() map[int]string {
	s := make(map[int]string)
	s[1] = `{"id":1,"first_name":"John","last_name":"Doe","age":37}`
	s[2] = `{"id":2,"first_name":"Abe","last_name":"Lincoln","age":52}`
	s[3] = `{"id":3,"first_name":"Bill","last_name":"Shakespeare","age":18}`
	s[4] = `{"id":4,"first_name":"Helen","last_name":"Keller","age":25}`

	return s
}
