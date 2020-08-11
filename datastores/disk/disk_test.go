package disk

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"testing"
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

			tf := newTestTableFile(t)
			defer tf.close()

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

			wantErr := ErrNoTable
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
				t.Fatal("File already exists for dsk.CreateTable test!!!")

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
	}

	for i, fn := range tests {
		testSetup(t)
		t.Run(strconv.Itoa(i), fn)
		testTeardown(t)
	}
}
