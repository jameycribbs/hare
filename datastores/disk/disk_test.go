package disk

import (
	"errors"
	"os"
	"reflect"
	"sort"
	"strconv"
	"testing"

	"github.com/jameycribbs/hare/dberr"
)

func TestAllDiskTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//New...

			dsk := newTestDisk(t)
			defer dsk.Close()

			want := "./testdata"
			got := dsk.path
			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}

			want = ".json"
			got = dsk.ext
			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}

			// DO I STILL NEED THESE?
			//tf := newTestTableFile(t)
			//defer tf.close()

			wantOffsets := make(map[int]int64)
			wantOffsets[1] = 0
			wantOffsets[2] = 101
			wantOffsets[3] = 160
			wantOffsets[4] = 224

			gotOffsets := dsk.tableFiles["contacts"].offsets

			if !reflect.DeepEqual(wantOffsets, gotOffsets) {
				t.Errorf("want %v; got %v", wantOffsets, gotOffsets)
			}
		},
		func(t *testing.T) {
			//Close...

			dsk := newTestDisk(t)
			dsk.Close()

			wantErr := dberr.ErrNoTable
			_, gotErr := dsk.ReadRec("contacts", 3)

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}

			got := dsk.tableFiles

			if nil != got {
				t.Errorf("want %v; got %v", nil, got)
			}
		},
		func(t *testing.T) {
			//CreateTable...

			if _, err := os.Stat("./testdata/newtable.json"); err == nil {
				t.Fatal("file already exists for dsk.CreateTable test")

			}

			dsk := newTestDisk(t)
			defer dsk.Close()

			err := dsk.CreateTable("newtable")
			if err != nil {
				t.Fatal(err)
			}

			if _, err = os.Stat("./testdata/newtable.json"); err != nil {
				t.Errorf("want %v; got %v", nil, err)
			}

			want := true
			got := dsk.TableExists("newtable")

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//CreateTable (TableExists error)...

			dsk := newTestDisk(t)
			defer dsk.Close()

			wantErr := dberr.ErrTableExists
			gotErr := dsk.CreateTable("contacts")

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//DeleteRec...

			dsk := newTestDisk(t)
			defer dsk.Close()

			err := dsk.DeleteRec("contacts", 3)
			if err != nil {
				t.Fatal(err)
			}

			want := dberr.ErrNoRecord
			_, got := dsk.ReadRec("contacts", 3)

			if !errors.Is(got, want) {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//DeleteRec (NoTable error)...

			dsk := newTestDisk(t)
			defer dsk.Close()

			wantErr := dberr.ErrNoTable
			gotErr := dsk.DeleteRec("nonexistent", 3)

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//GetLastID...

			dsk := newTestDisk(t)
			defer dsk.Close()

			want := 4
			got, err := dsk.GetLastID("contacts")
			if err != nil {
				t.Fatal(err)
			}

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//GetLastID (NoTable error)...

			dsk := newTestDisk(t)
			defer dsk.Close()

			wantErr := dberr.ErrNoTable
			_, gotErr := dsk.GetLastID("nonexistent")

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//IDs...

			dsk := newTestDisk(t)
			defer dsk.Close()

			want := []int{1, 2, 3, 4}
			got, err := dsk.IDs("contacts")
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

			dsk := newTestDisk(t)
			defer dsk.Close()

			wantErr := dberr.ErrNoTable
			_, gotErr := dsk.IDs("nonexistent")

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//InsertRec...

			dsk := newTestDisk(t)
			defer dsk.Close()

			err := dsk.InsertRec("contacts", 5, []byte(`{"id":5,"first_name":"Rex","last_name":"Stout","age":77}`))
			if err != nil {
				t.Fatal(err)
			}

			rec, err := dsk.ReadRec("contacts", 5)
			if err != nil {
				t.Fatal(err)
			}

			want := "{\"id\":5,\"first_name\":\"Rex\",\"last_name\":\"Stout\",\"age\":77}\n"
			got := string(rec)

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//InsertRec (NoTable error)...

			dsk := newTestDisk(t)
			defer dsk.Close()

			wantErr := dberr.ErrNoTable
			gotErr := dsk.InsertRec("nonexistent", 5, []byte(`{"id":5,"first_name":"Rex","last_name":"Stout","age":77}`))

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//InsertRec (IDExists error)...

			dsk := newTestDisk(t)
			defer dsk.Close()

			wantErr := dberr.ErrIDExists
			gotErr := dsk.InsertRec("contacts", 3, []byte(`{"id":3,"first_name":"Rex","last_name":"Stout","age":77}`))
			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}

			rec, err := dsk.ReadRec("contacts", 3)
			if err != nil {
				t.Fatal(err)
			}

			want := "{\"id\":3,\"first_name\":\"Bill\",\"last_name\":\"Shakespeare\",\"age\":18}\n"
			got := string(rec)

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//ReadRec...

			dsk := newTestDisk(t)
			defer dsk.Close()

			rec, err := dsk.ReadRec("contacts", 3)
			if err != nil {
				t.Fatal(err)
			}

			want := "{\"id\":3,\"first_name\":\"Bill\",\"last_name\":\"Shakespeare\",\"age\":18}\n"
			got := string(rec)

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//ReadRec (NoTable error)...

			dsk := newTestDisk(t)
			defer dsk.Close()

			wantErr := dberr.ErrNoTable
			_, gotErr := dsk.ReadRec("nonexistent", 3)

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//RemoveTable...

			if _, err := os.Stat("./testdata/contacts.json"); err != nil {
				t.Fatal(err)

			}

			dsk := newTestDisk(t)
			defer dsk.Close()

			err := dsk.RemoveTable("contacts")
			if err != nil {
				t.Fatal(err)
			}

			if _, err := os.Stat("./testdata/contacts.json"); !os.IsNotExist(err) {
				t.Errorf("want %v; got %v", os.ErrNotExist, err)
			}

			want := false
			got := dsk.TableExists("contacts")

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//RemoveTable (NoTable error)...

			dsk := newTestDisk(t)
			defer dsk.Close()

			wantErr := dberr.ErrNoTable
			gotErr := dsk.RemoveTable("nonexistent")

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//TableExists...

			dsk := newTestDisk(t)
			defer dsk.Close()

			want := true
			got := dsk.TableExists("contacts")

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}

			want = false
			got = dsk.TableExists("nonexistant")

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//TableNames...

			dsk := newTestDisk(t)
			defer dsk.Close()

			want := []string{"contacts"}
			got := dsk.TableNames()

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
		func(t *testing.T) {
			//UpdateRec...

			dsk := newTestDisk(t)
			defer dsk.Close()

			err := dsk.UpdateRec("contacts", 3, []byte(`{"id":3,"first_name":"William","last_name":"Shakespeare","age":77}`))
			if err != nil {
				t.Fatal(err)
			}

			rec, err := dsk.ReadRec("contacts", 3)
			if err != nil {
				t.Fatal(err)
			}

			want := "{\"id\":3,\"first_name\":\"William\",\"last_name\":\"Shakespeare\",\"age\":77}\n"
			got := string(rec)

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//UpdateRec (NoTable error)...

			dsk := newTestDisk(t)
			defer dsk.Close()

			wantErr := dberr.ErrNoTable
			gotErr := dsk.UpdateRec("nonexistent", 3, []byte(`{"id":3,"first_name":"William","last_name":"Shakespeare","age":77}`))

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//closeTable...

			dsk := newTestDisk(t)
			defer dsk.Close()

			err := dsk.closeTable("contacts")
			if err != nil {
				t.Errorf("want %v; got %v", nil, err)
			}
		},
		func(t *testing.T) {
			//closeTable (NoTable error)...

			dsk := newTestDisk(t)
			defer dsk.Close()

			wantErr := dberr.ErrNoTable
			gotErr := dsk.closeTable("nonexistent")

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
	}

	for i, fn := range tests {
		testSetup(t)
		t.Run(strconv.Itoa(i), fn)
		testTeardown(t)
	}
}
