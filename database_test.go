package hare

import (
	"fmt"
	"sort"
	"strconv"
	"testing"

	"github.com/jameycribbs/hare/datastores/disk"
	"github.com/jameycribbs/hare/datastores/ram"
	"github.com/jameycribbs/hare/dberr"
)

func TestCloseDatabaseTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//Close Disk Database...

			ds, err := disk.New("./testdata", ".json")
			if err != nil {
				t.Fatal(err)
			}

			db, err := New(ds)
			if err != nil {
				t.Fatal(err)
			}
			db.Close()

			checkErr(t, dberr.ErrNoTable, db.Find("contacts", 3, &Contact{}))

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
			//Close Ram Database...

			r, err := ram.New(seedData())
			if err != nil {
				t.Fatal(err)
			}

			db, err := New(r)
			if err != nil {
				t.Fatal(err)
			}
			db.Close()

			checkErr(t, dberr.ErrNoTable, db.Find("contacts", 3, &Contact{}))

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
	}

	for i, fn := range tests {
		t.Run(strconv.Itoa(i), fn)
	}
}

func TestNonMutatingDatabaseTests(t *testing.T) {
	var tests = []func(*Database) func(*testing.T){
		func(db *Database) func(*testing.T) {
			//New...

			return func(t *testing.T) {
				want := 4
				got := db.lastIDs["contacts"]
				if want != got {
					t.Errorf("want %v; got %v", want, got)
				}
			}
		},
		func(db *Database) func(*testing.T) {
			//TableExists...

			return func(t *testing.T) {
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
			}
		},
		func(db *Database) func(*testing.T) {
			//tableExists...

			return func(t *testing.T) {
				want := true
				got := db.tableExists("contacts")
				if want != got {
					t.Errorf("want %v; got %v", want, got)
				}

				want = false
				got = db.tableExists("nonexistent")
				if want != got {
					t.Errorf("want %v; got %v", want, got)
				}
			}
		},
	}

	runTestFns(t, tests)
}

func TestMutatingDatabaseTests(t *testing.T) {
	var tests = []func(*Database) func(*testing.T){
		func(db *Database) func(*testing.T) {
			//CreateTable...

			return func(t *testing.T) {
				err := db.CreateTable("newtable")
				if err != nil {
					t.Fatal(err)
				}

				want := true
				got := db.TableExists("newtable")

				if want != got {
					t.Errorf("want %v; got %v", want, got)
				}
			}
		},
		func(db *Database) func(*testing.T) {
			//CreateTable (TableExists error)...

			return func(t *testing.T) {
				checkErr(t, dberr.ErrTableExists, db.CreateTable("contacts"))
			}
		},
		func(db *Database) func(*testing.T) {
			//DropTable...

			return func(t *testing.T) {
				err := db.DropTable("contacts")
				if err != nil {
					t.Fatal(err)
				}

				want := false
				got := db.TableExists("contacts")

				if want != got {
					t.Errorf("want %v; got %v", want, got)
				}
			}
		},
		func(db *Database) func(*testing.T) {
			//DropTable (NoTable error)...

			return func(t *testing.T) {
				checkErr(t, dberr.ErrNoTable, db.DropTable("nonexistent"))
			}
		},
		func(db *Database) func(*testing.T) {
			//incrementLastID...

			return func(t *testing.T) {
				want := 5
				got := db.incrementLastID("contacts")

				if want != got {
					t.Errorf("want %v; got %v", want, got)
				}
			}
		},
	}

	runTestFns(t, tests)
}

func TestTableTests(t *testing.T) {
	var tests = []func(*Database) func(*testing.T){
		func(db *Database) func(*testing.T) {
			//IDs()...

			return func(t *testing.T) {
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
			}
		},
		func(db *Database) func(*testing.T) {
			//IDs() (NoTable error)...

			return func(t *testing.T) {
				_, gotErr := db.IDs("nonexistent")

				checkErr(t, dberr.ErrNoTable, gotErr)
			}
		},
	}

	runTestFns(t, tests)
}

func TestRecordTests(t *testing.T) {
	var tests = []func(*Database) func(*testing.T){
		func(db *Database) func(*testing.T) {
			//Delete...

			return func(t *testing.T) {
				err := db.Delete("contacts", 3)
				if err != nil {
					t.Fatal(err)
				}

				checkErr(t, dberr.ErrNoRecord, db.Find("contacts", 3, &Contact{}))
			}
		},
		func(db *Database) func(*testing.T) {
			//Delete (ErrNoTable error)...

			return func(t *testing.T) {
				checkErr(t, dberr.ErrNoTable, db.Delete("nonexistent", 3))
			}
		},
		func(db *Database) func(*testing.T) {
			//Delete (ErrNoRecord error)...

			return func(t *testing.T) {
				checkErr(t, dberr.ErrNoRecord, db.Delete("contacts", 99))
			}
		},
		func(db *Database) func(*testing.T) {
			//Find...

			return func(t *testing.T) {
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
			}
		},
		func(db *Database) func(*testing.T) {
			//Find (ErrNoRecord error)...

			return func(t *testing.T) {
				checkErr(t, dberr.ErrNoRecord, db.Find("contacts", 99, &Contact{}))
			}
		},
		func(db *Database) func(*testing.T) {
			//Insert...

			return func(t *testing.T) {
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
			}
		},
		func(db *Database) func(*testing.T) {
			//Insert (ErrNoTable error)...

			return func(t *testing.T) {
				_, gotErr := db.Insert("nonexistent", &Contact{FirstName: "Robin", LastName: "Williams", Age: 88})

				checkErr(t, dberr.ErrNoTable, gotErr)
			}
		},
		func(db *Database) func(*testing.T) {
			//Update...

			return func(t *testing.T) {
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
			}
		},
		func(db *Database) func(*testing.T) {
			//Update (ErrNoTable error)...

			return func(t *testing.T) {
				gotErr := db.Update("nonexistent", &Contact{ID: 4, FirstName: "Hazel", LastName: "Koller", Age: 26})

				checkErr(t, dberr.ErrNoTable, gotErr)
			}
		},
	}

	runTestFns(t, tests)
}
