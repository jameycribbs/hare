package hare

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"testing"

	"github.com/jameycribbs/hare/datastores/disk"
	"github.com/jameycribbs/hare/hare_err"
)

func TestAllDatabaseDiskTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//New...

			ds, err := disk.New("./testdata", ".json")
			if err != nil {
				t.Fatal(err)
			}

			db, err := New(ds)
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

			db := newTestDatabaseDisk(t)
			db.Close()

			wantErr := hare_err.NoTable
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

			db := newTestDatabaseDisk(t)
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
			//CreateTable (ErrTableExists)...

			db := newTestDatabaseDisk(t)
			defer db.Close()

			wantErr := hare_err.TableExists
			gotErr := db.CreateTable("contacts")

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//Delete...

			db := newTestDatabaseDisk(t)
			defer db.Close()

			err := db.Delete("contacts", 3)
			if err != nil {
				t.Fatal(err)
			}

			wantErr := hare_err.NoRecord
			gotErr := db.Find("contacts", 3, &Contact{})

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//Delete (ErrNoTable)...

			db := newTestDatabaseDisk(t)
			defer db.Close()

			wantErr := hare_err.NoTable
			gotErr := db.Delete("nonexistent", 3)

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//DropTable...

			db := newTestDatabaseDisk(t)
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
			//DropTable (ErrNoTable)...

			db := newTestDatabaseDisk(t)
			defer db.Close()

			wantErr := hare_err.NoTable
			gotErr := db.DropTable("nonexistent")

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//Find...

			db := newTestDatabaseDisk(t)
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
			//Find (ErrNoRecord)...

			db := newTestDatabaseDisk(t)
			defer db.Close()

			wantErr := hare_err.NoRecord
			gotErr := db.Find("contacts", 5, &Contact{})

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//IDs()...

			db := newTestDatabaseDisk(t)
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
			//IDs() (ErrNoTable)...

			db := newTestDatabaseDisk(t)
			defer db.Close()

			wantErr := hare_err.NoTable
			_, gotErr := db.IDs("nonexistent")

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//Insert...

			db := newTestDatabaseDisk(t)
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
			//Insert (ErrNoTable)...

			db := newTestDatabaseDisk(t)
			defer db.Close()

			wantErr := hare_err.NoTable
			_, gotErr := db.Insert("nonexistent", &Contact{FirstName: "Robin", LastName: "Williams", Age: 88})

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//TableExists...

			db := newTestDatabaseDisk(t)
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

			db := newTestDatabaseDisk(t)
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
			//Update (ErrNoTable)...

			db := newTestDatabaseDisk(t)
			defer db.Close()

			wantErr := hare_err.NoTable
			gotErr := db.Update("nonexistent", &Contact{ID: 4, FirstName: "Hazel", LastName: "Koller", Age: 26})

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//incrementLastID...

			db := newTestDatabaseDisk(t)
			defer db.Close()

			want := 5
			got := db.incrementLastID("contacts")

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//tableExists...

			db := newTestDatabaseDisk(t)
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
		testSetup(t)
		t.Run(strconv.Itoa(i), fn)
		testTeardown(t)
	}
}

func newTestDatabaseDisk(t *testing.T) *Database {
	ds, err := disk.New("./testdata", ".json")
	if err != nil {
		t.Fatal(err)
	}

	db, err := New(ds)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func testSetup(t *testing.T) {
	testRemoveFiles(t)

	cmd := exec.Command("cp", "./testdata/contacts.bak", "./testdata/contacts.json")
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
}

func testTeardown(t *testing.T) {
	testRemoveFiles(t)
}

func testRemoveFiles(t *testing.T) {
	filesToRemove := []string{"contacts.json", "newtable.json"}

	for _, f := range filesToRemove {
		err := os.Remove("./testdata/" + f)
		if err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}
	}
}
