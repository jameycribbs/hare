package ram

import (
	"errors"
	"reflect"
	"sort"
	"testing"

	"github.com/jameycribbs/hare/dberr"
)

func TestNewCloseRamTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//New...

			ram := newTestRam(t)
			defer ram.Close()

			want := make(map[int][]byte)
			want[1] = []byte(`{"id":1,"first_name":"John","last_name":"Doe","age":37}`)
			want[2] = []byte(`{"id":2,"first_name":"Abe","last_name":"Lincoln","age":52}`)
			want[3] = []byte(`{"id":3,"first_name":"Bill","last_name":"Shakespeare","age":18}`)
			want[4] = []byte(`{"id":4,"first_name":"Helen","last_name":"Keller","age":25}`)

			got := ram.tables["contacts"].records

			if !reflect.DeepEqual(want, got) {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//Close...

			ram := newTestRam(t)
			ram.Close()

			wantErr := dberr.ErrNoTable
			_, gotErr := ram.ReadRec("contacts", 3)

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}

			got := ram.tables

			if nil != got {
				t.Errorf("want %v; got %v", nil, got)
			}
		},
	}

	runTestFns(t, tests)
}

func TestCreateTableRamTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//CreateTable...

			ram := newTestRam(t)
			defer ram.Close()

			err := ram.CreateTable("newtable")
			if err != nil {
				t.Fatal(err)
			}

			want := true
			got := ram.TableExists("newtable")

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//CreateTable (TableExists error)...

			ram := newTestRam(t)
			defer ram.Close()

			wantErr := dberr.ErrTableExists
			gotErr := ram.CreateTable("contacts")

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
	}

	runTestFns(t, tests)
}

func TestDeleteRecRamTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//DeleteRec...

			ram := newTestRam(t)
			defer ram.Close()

			err := ram.DeleteRec("contacts", 3)
			if err != nil {
				t.Fatal(err)
			}

			want := dberr.ErrNoRecord
			_, got := ram.ReadRec("contacts", 3)

			if !errors.Is(got, want) {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//DeleteRec (NoTable error)...

			ram := newTestRam(t)
			defer ram.Close()

			wantErr := dberr.ErrNoTable
			gotErr := ram.DeleteRec("nonexistent", 3)

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
	}

	runTestFns(t, tests)
}

func TestGetLastIDRamTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//GetLastID...

			ram := newTestRam(t)
			defer ram.Close()

			want := 4
			got, err := ram.GetLastID("contacts")
			if err != nil {
				t.Fatal(err)
			}

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//GetLastID (NoTable error)...

			ram := newTestRam(t)
			defer ram.Close()

			wantErr := dberr.ErrNoTable
			_, gotErr := ram.GetLastID("nonexistent")

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
	}

	runTestFns(t, tests)
}

func TestIDsRamTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//IDs...

			ram := newTestRam(t)
			defer ram.Close()

			want := []int{1, 2, 3, 4}
			got, err := ram.IDs("contacts")
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
			//IDs (NoTable error)...

			ram := newTestRam(t)
			defer ram.Close()

			wantErr := dberr.ErrNoTable
			_, gotErr := ram.IDs("nonexistent")

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
	}

	runTestFns(t, tests)
}

func TestInsertRecRamTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//InsertRec...

			ram := newTestRam(t)
			defer ram.Close()

			err := ram.InsertRec("contacts", 5, []byte(`{"id":5,"first_name":"Rex","last_name":"Stout","age":77}`))
			if err != nil {
				t.Fatal(err)
			}

			rec, err := ram.ReadRec("contacts", 5)
			if err != nil {
				t.Fatal(err)
			}

			want := `{"id":5,"first_name":"Rex","last_name":"Stout","age":77}`
			got := string(rec)

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//InsertRec (NoTable error)...

			ram := newTestRam(t)
			defer ram.Close()

			wantErr := dberr.ErrNoTable
			gotErr := ram.InsertRec("nonexistent", 5, []byte(`{"id":5,"first_name":"Rex","last_name":"Stout","age":77}`))

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//InsertRec (IDExists error)...

			ram := newTestRam(t)
			defer ram.Close()

			wantErr := dberr.ErrIDExists
			gotErr := ram.InsertRec("contacts", 3, []byte(`{"id":3,"first_name":"Rex","last_name":"Stout","age":77}`))
			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}

			rec, err := ram.ReadRec("contacts", 3)
			if err != nil {
				t.Fatal(err)
			}

			want := `{"id":3,"first_name":"Bill","last_name":"Shakespeare","age":18}`
			got := string(rec)

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
	}

	runTestFns(t, tests)
}

func TestRecRecRamTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//ReadRec...

			ram := newTestRam(t)
			defer ram.Close()

			rec, err := ram.ReadRec("contacts", 3)
			if err != nil {
				t.Fatal(err)
			}

			want := `{"id":3,"first_name":"Bill","last_name":"Shakespeare","age":18}`
			got := string(rec)

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//ReadRec (NoTable error)...

			ram := newTestRam(t)
			defer ram.Close()

			wantErr := dberr.ErrNoTable
			_, gotErr := ram.ReadRec("nonexistent", 3)

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
	}

	runTestFns(t, tests)
}

func TestRemoveTableRamTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//RemoveTable...

			ram := newTestRam(t)
			defer ram.Close()

			err := ram.RemoveTable("contacts")
			if err != nil {
				t.Fatal(err)
			}

			want := false
			got := ram.TableExists("contacts")

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//RemoveTable (NoTable error)...

			ram := newTestRam(t)
			defer ram.Close()

			wantErr := dberr.ErrNoTable
			gotErr := ram.RemoveTable("nonexistent")

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
	}

	runTestFns(t, tests)
}

func TestTableExistsRamTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//TableExists...

			ram := newTestRam(t)
			defer ram.Close()

			want := true
			got := ram.TableExists("contacts")

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}

			want = false
			got = ram.TableExists("nonexistant")

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
	}

	runTestFns(t, tests)
}

func TestTableNamesRamTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//TableNames...

			ram := newTestRam(t)
			defer ram.Close()

			want := []string{"contacts"}
			got := ram.TableNames()

			sort.Strings(got)

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
	}

	runTestFns(t, tests)
}

func TestUpdateRecRamTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//UpdateRec...

			ram := newTestRam(t)
			defer ram.Close()

			err := ram.UpdateRec("contacts", 3, []byte(`{"id":3,"first_name":"William","last_name":"Shakespeare","age":77}`))
			if err != nil {
				t.Fatal(err)
			}

			rec, err := ram.ReadRec("contacts", 3)
			if err != nil {
				t.Fatal(err)
			}

			want := `{"id":3,"first_name":"William","last_name":"Shakespeare","age":77}`
			got := string(rec)

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//UpdateRec (NoTable error)...

			ram := newTestRam(t)
			defer ram.Close()

			wantErr := dberr.ErrNoTable
			gotErr := ram.UpdateRec("nonexistent", 3, []byte(`{"id":3,"first_name":"William","last_name":"Shakespeare","age":77}`))

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
	}

	runTestFns(t, tests)
}
