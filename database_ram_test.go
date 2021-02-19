package hare

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"testing"

	"github.com/jameycribbs/hare/datastores/ram"
	"github.com/jameycribbs/hare/dberr"
)

//gocyclo:ignore
func TestAllDatabaseRamTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//New...

			r, err := ram.New(seedData())
			if err != nil {
				t.Fatal(err)
			}

			db, err := New(r)
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()

			want := 4
			got := db.lastIDs["contacts"]
			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//Close...

			db := newTestDatabaseRam(t)
			db.Close()

			wantErr := dberr.ErrNoTable
			gotErr := db.Find("contacts", 3, &Contact{})

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}

			gotStore := db.store
			if nil != gotStore {
				t.Errorf("want %v; got %v", nil, gotStore)
			}

			gotLocks := db.locks
			if nil != gotLocks {
				t.Errorf("want %v; got %v", nil, gotLocks)
			}

			gotLastIDs := db.lastIDs
			if nil != gotLastIDs {
				t.Errorf("want %v; got %v", nil, gotLastIDs)
			}
		},
		func(t *testing.T) {
			//CreateTable...

			db := newTestDatabaseRam(t)
			defer db.Close()

			err := db.CreateTable("newtable")
			if err != nil {
				t.Fatal(err)
			}

			want := true
			got := db.TableExists("newtable")

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//CreateTable (TableExists error)...

			db := newTestDatabaseRam(t)
			defer db.Close()

			wantErr := dberr.ErrTableExists
			gotErr := db.CreateTable("contacts")

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//Delete...

			db := newTestDatabaseRam(t)
			defer db.Close()

			err := db.Delete("contacts", 3)
			if err != nil {
				t.Fatal(err)
			}

			wantErr := dberr.ErrNoRecord
			gotErr := db.Find("contacts", 3, &Contact{})

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//Delete (NoTable error)...

			db := newTestDatabaseRam(t)
			defer db.Close()

			wantErr := dberr.ErrNoTable
			gotErr := db.Delete("nonexistent", 3)

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//DropTable...

			db := newTestDatabaseRam(t)
			defer db.Close()

			err := db.DropTable("contacts")
			if err != nil {
				t.Fatal(err)
			}

			want := false
			got := db.TableExists("contacts")

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//DropTable (NoTable error)...

			db := newTestDatabaseRam(t)
			defer db.Close()

			wantErr := dberr.ErrNoTable
			gotErr := db.DropTable("nonexistent")

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//Find...

			db := newTestDatabaseRam(t)
			defer db.Close()

			c := Contact{}

			err := db.Find("contacts", 2, &c)
			if err != nil {
				t.Fatal(err)
			}

			want := "Abe Lincoln is 52"
			got := fmt.Sprintf("%s %s is %d", c.FirstName, c.LastName, c.Age)

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//Find (NoRecord error)...

			db := newTestDatabaseRam(t)
			defer db.Close()

			wantErr := dberr.ErrNoRecord
			gotErr := db.Find("contacts", 5, &Contact{})

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//IDs()...

			db := newTestDatabaseRam(t)
			defer db.Close()

			want := []int{1, 2, 3, 4}
			got, err := db.IDs("contacts")
			if err != nil {
				t.Fatal(err)
			}

			sort.Ints(got)

			if len(want) != len(got) {
				t.Errorf("want %v; got %v", want, got)
			} else {

				for i := range want {
					if want[i] != got[i] {
						t.Errorf("want %v; got %v", want, got)
					}
				}
			}
		},
		func(t *testing.T) {
			//IDs() (NoTable error)...

			db := newTestDatabaseRam(t)
			defer db.Close()

			wantErr := dberr.ErrNoTable
			_, gotErr := db.IDs("nonexistent")

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//Insert...

			db := newTestDatabaseRam(t)
			defer db.Close()

			wantInt := 5
			gotInt, err := db.Insert("contacts", &Contact{FirstName: "Robin", LastName: "Williams", Age: 88})
			if err != nil {
				t.Fatal(err)
			}

			if wantInt != gotInt {
				t.Errorf("want %v; got %v", wantInt, gotInt)
			}

			c := Contact{}

			err = db.Find("contacts", 5, &c)
			if err != nil {
				t.Fatal(err)
			}

			want := "Robin Williams is 88"
			got := fmt.Sprintf("%s %s is %d", c.FirstName, c.LastName, c.Age)

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//Insert (NoTable error)...

			db := newTestDatabaseRam(t)
			defer db.Close()

			wantErr := dberr.ErrNoTable
			_, gotErr := db.Insert("nonexistent", &Contact{FirstName: "Robin", LastName: "Williams", Age: 88})

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//TableExists...

			db := newTestDatabaseRam(t)
			defer db.Close()

			want := true
			got := db.TableExists("contacts")

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}

			want = false
			got = db.TableExists("nonexistent")

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//Update...

			db := newTestDatabaseRam(t)
			defer db.Close()

			err := db.Update("contacts", &Contact{ID: 4, FirstName: "Hazel", LastName: "Koller", Age: 26})
			if err != nil {
				t.Fatal(err)
			}

			c := Contact{}

			err = db.Find("contacts", 4, &c)
			if err != nil {
				t.Fatal(err)
			}

			want := "Hazel Koller is 26"
			got := fmt.Sprintf("%s %s is %d", c.FirstName, c.LastName, c.Age)

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//Update (NoTable error)...

			db := newTestDatabaseRam(t)
			defer db.Close()

			wantErr := dberr.ErrNoTable
			gotErr := db.Update("nonexistent", &Contact{ID: 4, FirstName: "Hazel", LastName: "Koller", Age: 26})

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//incrementLastID...

			db := newTestDatabaseRam(t)
			defer db.Close()

			want := 5
			got := db.incrementLastID("contacts")

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//tableExists...

			db := newTestDatabaseRam(t)
			defer db.Close()

			want := true
			got := db.tableExists("contacts")

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}

			want = false
			got = db.TableExists("nonexistent")

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
	}

	for i, fn := range tests {
		t.Run(strconv.Itoa(i), fn)
	}
}

func newTestDatabaseRam(t *testing.T) *Database {
	r, err := ram.New(seedData())
	if err != nil {
		t.Fatal(err)
	}

	db, err := New(r)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func seedData() map[string]map[int]string {
	tblMap := make(map[string]map[int]string)
	contactsMap := make(map[int]string)

	contactsMap[1] = `{"id":1,"first_name":"John","last_name":"Doe","age":37}`
	contactsMap[2] = `{"id":2,"first_name":"Abe","last_name":"Lincoln","age":52}`
	contactsMap[3] = `{"id":3,"first_name":"Bill","last_name":"Shakespeare","age":18}`
	contactsMap[4] = `{"id":4,"first_name":"Helen","last_name":"Keller","age":25}`

	tblMap["contacts"] = contactsMap

	return tblMap
}
