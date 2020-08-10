package disk

import (
	"bufio"
	"errors"
	"reflect"
	"sort"
	"strconv"
	"testing"
)

func TestAllTableFileTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//New...

			tf := newTestTableFile(t)
			defer tf.close()

			want := make(map[int]int64)
			want[1] = 0
			want[2] = 101
			want[3] = 160
			want[4] = 224

			got := tf.offsets

			if !reflect.DeepEqual(want, got) {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//close...

			tf := newTestTableFile(t)
			tf.close()

			_, gotErr := tf.readRec(3)
			if !errors.Is(gotErr, ErrNoRecord) {
				t.Errorf("want %v; got %v", ErrNoRecord, gotErr)
			}

			got := tf.offsets

			if nil != got {
				t.Errorf("want %v; got %v", nil, got)
			}
		},
		func(t *testing.T) {
			//deleteRec...

			tf := newTestTableFile(t)
			defer tf.close()

			offset := tf.offsets[3]

			err := tf.deleteRec(3)
			if err != nil {
				t.Fatal(err)
			}

			want := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\n"

			r := bufio.NewReader(tf.ptr)

			if _, err := tf.ptr.Seek(offset, 0); err != nil {
				t.Fatal(err)
			}

			rec, err := r.ReadBytes('\n')
			if err != nil {
				t.Fatal(err)
			}
			got := string(rec)

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//getLastID...

			tf := newTestTableFile(t)
			defer tf.close()

			want := 4
			got := tf.getLastID()

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//ids...

			tf := newTestTableFile(t)
			defer tf.close()

			want := []int{1, 2, 3, 4}
			got := tf.ids()
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
			//offsetForWritingRec...

			tf := newTestTableFile(t)
			defer tf.close()

			tests := []struct {
				recLen int
				want   int
			}{
				{45, 284},
				{44, 56},
			}

			for _, tt := range tests {
				want := int64(tt.want)
				got, err := tf.offsetForWritingRec(tt.recLen)
				if err != nil {
					t.Fatal(err)
				}
				if want != got {
					t.Errorf("want %v; got %v", want, got)
				}
			}
		},
		func(t *testing.T) {
			//offsetToFitRec...

			tf := newTestTableFile(t)
			defer tf.close()

			tests := []struct {
				recLen  int
				want    int
				wanterr error
			}{
				{284, 0, dummiesTooShortError{}},
				{44, 56, nil},
			}

			for _, tt := range tests {
				want := int64(tt.want)
				got, goterr := tf.offsetToFitRec(tt.recLen)
				if !((want == got) && (errors.Is(goterr, tt.wanterr))) {
					t.Errorf("want %v; wanterr %v; got %v; goterr %v", want, tt.wanterr, got, goterr)
				}
			}
		},
		func(t *testing.T) {
			//overwriteRec...

			tf := newTestTableFile(t)
			defer tf.close()

			offset := tf.offsets[3]

			err := tf.overwriteRec(160, 64)
			if err != nil {
				t.Fatal(err)
			}

			want := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\n"

			r := bufio.NewReader(tf.ptr)

			if _, err := tf.ptr.Seek(offset, 0); err != nil {
				t.Fatal(err)
			}

			rec, err := r.ReadBytes('\n')
			if err != nil {
				t.Fatal(err)
			}
			got := string(rec)

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//readRec...

			tf := newTestTableFile(t)
			defer tf.close()

			rec, err := tf.readRec(3)
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
			//updateRec (fits on same line)...

			tf := newTestTableFile(t)
			defer tf.close()

			err := tf.updateRec(3, []byte("{\"id\":3,\"first_name\":\"Bill\",\"last_name\":\"Shakespeare\",\"age\":92}"))
			if err != nil {
				t.Fatal(err)
			}

			wantOffset := int64(160)
			gotOffset := tf.offsets[3]

			if wantOffset != gotOffset {
				t.Errorf("want %v; got %v", wantOffset, gotOffset)
			}

			rec, err := tf.readRec(3)
			if err != nil {
				t.Fatal(err)
			}

			want := "{\"id\":3,\"first_name\":\"Bill\",\"last_name\":\"Shakespeare\",\"age\":92}\n"
			got := string(rec)

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//updateRec (does not fit on same line)...

			tf := newTestTableFile(t)
			defer tf.close()

			err := tf.updateRec(3, []byte("{\"id\":3,\"first_name\":\"William\",\"last_name\":\"Shakespeare\",\"age\":18}"))
			if err != nil {
				t.Fatal(err)
			}

			wantOffset := int64(284)
			gotOffset := tf.offsets[3]

			if wantOffset != gotOffset {
				t.Errorf("want %v; got %v", wantOffset, gotOffset)
			}

			rec, err := tf.readRec(3)
			if err != nil {
				t.Fatal(err)
			}

			want := "{\"id\":3,\"first_name\":\"William\",\"last_name\":\"Shakespeare\",\"age\":18}\n"
			got := string(rec)

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//padRec...

			want := "\nXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
			got := string(padRec(50))

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
