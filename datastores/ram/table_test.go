package ram

import (
	"errors"
	"reflect"
	"sort"
	"testing"

	"github.com/jameycribbs/hare/dberr"
)

func TestNewTableTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//New...

			tbl := newTestTable(t)

			want := make(map[int][]byte)
			want[1] = []byte(`{"id":1,"first_name":"John","last_name":"Doe","age":37}`)
			want[2] = []byte(`{"id":2,"first_name":"Abe","last_name":"Lincoln","age":52}`)
			want[3] = []byte(`{"id":3,"first_name":"Bill","last_name":"Shakespeare","age":18}`)
			want[4] = []byte(`{"id":4,"first_name":"Helen","last_name":"Keller","age":25}`)

			got := tbl.records

			if !reflect.DeepEqual(want, got) {
				t.Errorf("want %v; got %v", want, got)
			}
		},
	}

	runTestFns(t, tests)
}

func TestDeleteRecTableTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//deleteRec...

			tbl := newTestTable(t)

			err := tbl.deleteRec(3)
			if err != nil {
				t.Fatal(err)
			}

			wantErr := dberr.ErrNoRecord
			_, gotErr := tbl.readRec(3)
			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
	}

	runTestFns(t, tests)
}

func TestGetLastIDTableTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//getLastID...

			tbl := newTestTable(t)

			want := 4
			got := tbl.getLastID()

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
	}

	runTestFns(t, tests)
}

func TestIDsTableTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//ids...

			tbl := newTestTable(t)

			want := []int{1, 2, 3, 4}
			got := tbl.ids()
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
	}

	runTestFns(t, tests)
}

func TestReadRecTableTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//readRec...

			tbl := newTestTable(t)

			rec, err := tbl.readRec(3)
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

func TestWriteRecTableTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//writeRec

			tbl := newTestTable(t)

			want := []byte(`{"id":3,"first_name":"Bill","last_name":"Shakespeare","age":92}`)
			tbl.writeRec(3, want)

			got, err := tbl.readRec(3)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(want, got) {
				t.Errorf("want %v; got %v", want, got)
			}
		},
	}

	runTestFns(t, tests)
}
